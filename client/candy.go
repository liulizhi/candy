package candy

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/dearcode/candy/meta"
	"github.com/dearcode/candy/util"
	"github.com/dearcode/candy/util/log"
)

const (
	networkTimeout = time.Second * 3
)

// MessageHandler 接收服务器端推送来的消息
type MessageHandler interface {
	// OnRecv 这函数理论上是多线程调用，客户端需要注意下
	OnRecv(id int64, method int, group int64, from int64, to int64, body string)

	// OnError 连接被服务器断开，或其它错误
	OnError(msg string)

	// OnUnHealth 连接异常
	OnUnHealth(msg string)
}

type CandyClient struct {
	host    string
	stop    bool
	conn    *grpc.ClientConn
	api     meta.GateClient
	handler MessageHandler
	stream  meta.Gate_MessageStreamClient
	health  healthpb.HealthClient
}

func NewCandyClient(host string, handler MessageHandler) *CandyClient {
	return &CandyClient{host: host, handler: handler}
}

// Start 连接服务端.
func (c *CandyClient) Start() (err error) {
	if c.conn, err = grpc.Dial(c.host, grpc.WithInsecure(), grpc.WithTimeout(networkTimeout)); err != nil {
		return
	}

	c.api = meta.NewGateClient(c.conn)
	if c.stream, err = c.api.MessageStream(context.Background()); err != nil {
		return
	}

	go c.loopRecvMessage()

	//健康检查
	c.health = healthpb.NewHealthClient(c.conn)
	go c.healthCheck()

	return
}

// Stop 断开到服务器连接.
func (c *CandyClient) Stop() error {
	c.stop = true
	c.stream.CloseSend()
	return c.conn.Close()
}

func (c *CandyClient) Register(user, passwd string) (int64, error) {
	req := &meta.GateRegisterRequest{User: user, Password: passwd}
	resp, err := c.api.Register(context.Background(), req)
	if err != nil {
		return -1, err
	}

	return resp.ID, resp.Header.Error()
}

func (c *CandyClient) Login(user, passwd string) (int64, error) {
	req := &meta.GateUserLoginRequest{User: user, Password: passwd}
	resp, err := c.api.Login(context.Background(), req)
	if err != nil {
		return -1, err
	}

	return resp.ID, resp.Header.Error()
}

func (c *CandyClient) Logout(user string) error {
	req := &meta.GateUserLogoutRequest{User: user}
	resp, err := c.api.Logout(context.Background(), req)
	if err != nil {
		return err
	}

	return resp.Header.Error()
}

func (c *CandyClient) UpdateUserInfo(user, nickName string, avatar []byte) (int64, error) {
	req := &meta.GateUpdateUserInfoRequest{User: user, NickName: nickName, Avatar: avatar}
	resp, err := c.api.UpdateUserInfo(context.Background(), req)
	if err != nil {
		return -1, err
	}

	return resp.ID, resp.Header.Error()
}

func (c *CandyClient) UpdateUserPassword(user, passwd string) (int64, error) {
	req := &meta.GateUpdateUserPasswordRequest{User: user, Password: passwd}
	resp, err := c.api.UpdateUserPassword(context.Background(), req)
	if err != nil {
		return -1, err
	}

	return resp.ID, resp.Header.Error()
}

func (c *CandyClient) GetUserInfoByName(user string) ([]byte, error) {
	userInfo, err := c.getUserInfoByName(user)
	if err != nil {
		return nil, err
	}

	data, err := encodeJSON(userInfo)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *CandyClient) getUserInfoByName(user string) (*UserInfo, error) {
	req := &meta.GateGetUserInfoRequest{Type: 0, UserName: user}
	resp, err := c.api.GetUserInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{ID: resp.ID, Name: resp.User, NickName: resp.NickName, Avatar: resp.Avatar}
	return userInfo, resp.Header.Error()
}

func (c *CandyClient) GetUserInfoByID(userID int64) ([]byte, error) {
	userInfo, err := c.getUserInfoByID(userID)
	if err != nil {
		return nil, err
	}

	data, err := encodeJSON(userInfo)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *CandyClient) getUserInfoByID(userID int64) (*UserInfo, error) {
	req := &meta.GateGetUserInfoRequest{Type: 1, UserID: userID}
	resp, err := c.api.GetUserInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{ID: resp.ID, Name: resp.User, NickName: resp.NickName, Avatar: resp.Avatar}
	return userInfo, resp.Header.Error()
}

func (c *CandyClient) AddFriend(userID int64, confirm bool, msg string) (bool, error) {
	req := &meta.GateAddFriendRequest{UserID: userID, Confirm: confirm, Msg: msg}
	resp, err := c.api.AddFriend(context.Background(), req)
	if err != nil {
		return false, err
	}

	return resp.Confirm, resp.Header.Error()
}

func (c *CandyClient) LoadFriendList() ([]byte, error) {
	req := &meta.GateLoadFriendListRequest{}
	resp, err := c.api.LoadFriendList(context.Background(), req)
	if err != nil {
		return nil, err
	}

	friendList := &FriendList{Users: resp.Users}
	data, err := encodeJSON(friendList)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 支持模糊查询，返回对应用户的列表
func (c *CandyClient) FindUser(user string) ([]byte, error) {
	req := &meta.GateFindUserRequest{User: user}
	resp, err := c.api.FindUser(context.Background(), req)
	if err != nil {
		return nil, err
	}

	users := make([]*UserInfo, 0)
	for _, matchUser := range resp.Users {
		userInfo, err := c.getUserInfoByName(matchUser)
		if err != nil {
			return nil, err
		}
		users = append(users, userInfo)
	}
	userList := &UserList{Users: users}
	data, err := encodeJSON(userList)
	if err != nil {
		return nil, err
	}

	return data, resp.Header.Error()
}

func (c *CandyClient) FileExist(key string) (bool, error) {
	req := &meta.GateCheckFileRequest{Names: []string{key}}
	resp, err := c.api.CheckFile(context.Background(), req)
	if err != nil {
		return false, err
	}

	if err = resp.Header.Error(); err != nil {
		return false, err
	}

	if len(resp.Names) == 0 {
		return true, nil
	}

	return false, nil
}

func (c *CandyClient) FileUpload(data []byte) (string, error) {
	md5 := string(util.MD5(data))
	exist, err := c.FileExist(md5)
	if err != nil {
		return md5, err
	}
	//已有别人上传过了
	if exist {
		return md5, nil
	}

	req := &meta.GateUploadFileRequest{File: data}
	resp, err := c.api.UploadFile(context.Background(), req)
	if err != nil {
		return md5, err
	}

	return md5, resp.Header.Error()
}

func (c *CandyClient) FileDownload(key string) ([]byte, error) {
	req := &meta.GateDownloadFileRequest{Names: []string{key}}
	resp, err := c.api.DownloadFile(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.Files[key], resp.Header.Error()
}

// SendMessage 向服务器发送消息.
func (c *CandyClient) SendMessage(from, group, user int64, body string) error {
	msg := &meta.Message{From: from, Group: group, User: user, Body: body}
	if err := c.stream.Send(msg); err != nil {
		return err
	}
	return nil
}

// loopRecvMessage 一直接收服务器返回消息, 直到退出.
func (c *CandyClient) loopRecvMessage() {
	log.Debugf("start loopRecvMessage")

	for !c.stop {
		msg, err := c.stream.Recv()
		log.Debugf("recv:%v, err:%v", msg, err)
		if err != nil {
			// 这里不退出会死循环
			c.handler.OnError(err.Error())
			break
		}

		c.handler.OnRecv(msg.ID, int(msg.Method), msg.Group, msg.From, msg.User, msg.Body)
	}
}

func (c *CandyClient) healthCheck() {
	for !c.stop {
		time.Sleep(time.Second)
		req := &healthpb.HealthCheckRequest{
			Service: "",
		}

		_, err := c.health.Check(context.Background(), req)
		if err != nil {
			c.handler.OnUnHealth(err.Error())
		}
	}
}

// Heartbeat 向服务器发送心跳信息
func (c *CandyClient) Heartbeat() error {
	req := &meta.GateHeartbeatRequest{}
	_, err := c.api.Heartbeat(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

// CreateGroup 创建群组
func (c *CandyClient) CreateGroup(name string) (int64, error) {
	req := &meta.GateCreateGroupRequest{GroupName: name}
	resp, err := c.api.CreateGroup(context.Background(), req)
	if err != nil {
		return -1, err
	}

	return resp.ID, nil
}

// LoadGroupList 拉取群组列表
func (c *CandyClient) LoadGroupList() ([]byte, error) {
	req := &meta.GateLoadGroupListRequest{}
	resp, err := c.api.LoadGroupList(context.Background(), req)
	if err != nil {
		return nil, err
	}

	var groups []*GroupInfo
	for _, group := range resp.Groups {
		groups = append(groups, &GroupInfo{ID: group.ID, Name: group.Name, Users: group.Users})
	}

	groupList := &GroupList{Groups: groups}
	data, err := encodeJSON(groupList)
	if err != nil {
		return nil, err
	}

	return data, nil
}
