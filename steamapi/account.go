package steamapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bang-go/network/httpx"
	"github.com/bang-go/util"
	"net/http"
	"time"
)

const (
	// UrlGetAccountList 使用登录令牌获取游戏服务器帐户列表。
	UrlGetAccountList = "https://api.steampowered.com/IGameServersService/GetAccountList/v1/"
	// UrlCreateAccount 创建一个永久游戏服务器帐户。
	UrlCreateAccount = "https://api.steampowered.com/IGameServersService/CreateAccount/v1/"
	// UrlSetMemo 此方法改变与游戏服务器帐户关联的备注。 备注不会对帐户产生任何影响。 备注在 GetAccountList 响应中出现，只用于提示帐户用途。
	UrlSetMemo = "https://api.steampowered.com/IGameServersService/SetMemo/v1/"
	// UrlResetLoginToken 为指定游戏服务器生成新的登录令牌。
	UrlResetLoginToken = "https://api.steampowered.com/IGameServersService/ResetLoginToken/v1/"
	// UrlDeleteAccount 删除一个永久的游戏服务器帐户。
	UrlDeleteAccount = "https://api.steampowered.com/IGameServersService/DeleteAccount/v1/"
	// UrlGetAccountPublicInfo 获取给定游戏服务器帐户的公开信息。
	UrlGetAccountPublicInfo = "https://api.steampowered.com/IGameServersService/GetAccountPublicInfo/v1/"
	// UrlQueryLoginToken 查询指定令牌的状态，您必须拥有该令牌方可查询。
	UrlQueryLoginToken = "https://api.steampowered.com/IGameServersService/QueryLoginToken/v1/"
	// UrlGetServerSteamIDsByIP 根据 IP 列表，获取服务器 SteamID 列表。
	UrlGetServerSteamIDsByIP = "https://api.steampowered.com/IGameServersService/GetServerSteamIDsByIP/v1/"
	// UrlGetServerIPsBySteamID 根据 SteamID 列表，获取服务器 IP 地址列表。
	UrlGetServerIPsBySteamID = "https://api.steampowered.com/IGameServersService/GetServerIPsBySteamID/v1/"
)

const (
	DefaultReqTimeout = 10 * time.Second
)

type GameSteamServer struct {
	SteamId     string `json:"steamid"`
	AppId       int    `json:"appid"`
	LoginToken  string `json:"login_token"`
	IsDeleted   bool   `json:"is_deleted"`
	IsExpired   bool   `json:"is_expired"`
	RtLastLogon int64  `json:"rt_last_logon"`
}

type AccountListResp struct {
	Response struct {
		Servers        []GameSteamServer `json:"servers"`
		IsBanned       bool              `json:"is_banned"`
		Expires        int64             `json:"expires"`
		Actor          string            `json:"actor"`
		LastActionTime int64             `json:"last_action_time"`
	} `json:"response"`
}
type CreatAccountResp struct {
	Response struct {
		SteamId    string `json:"steamid"`
		LoginToken string `json:"login_token"`
	} `json:"response"`
}

type GetAccountPublicInfoResp struct {
	Response struct {
		SteamId string `json:"steamid"`
		AppId   int    `json:"appid"`
	} `json:"response"`
}

type ResetLoginTokenResp struct {
	Response struct {
		LoginToken string `json:"login_token"`
	} `json:"response"`
}

type QueryLoginTokenResp struct {
	Response struct {
		SteamId  string `json:"steamid"`
		IsBanned bool   `json:"is_banned"`
		Expires  int64  `json:"expires"`
	} `json:"response"`
}

type GameSteamServerAddr struct {
	SteamId string `json:"steamid"`
	Addr    string `json:"addr"`
}

type GetServerSteamIDsByIPResp struct {
	Response struct {
		Servers []GameSteamServerAddr `json:"servers"`
	} `json:"response"`
}

type GetServerIPsBySteamIDResp struct {
	Response struct {
		Servers []GameSteamServerAddr `json:"servers"`
	} `json:"response"`
}

type GameServersService interface {
	GetAccountList() (*AccountListResp, error)
	CreateAccount(appId int, memo string) (*CreatAccountResp, error)
	SetMemo(steamId string, memo string) error
	ResetLoginToken(steamId string) (*ResetLoginTokenResp, error)
	DeleteAccount(steamId string) error
	GetAccountPublicInfo(steamId string) (*GetAccountPublicInfoResp, error)
	QueryLoginToken(loginToken string) (*QueryLoginTokenResp, error)
	GetServerSteamIDsByIP(serverIps []string) (resp *GetServerSteamIDsByIPResp, err error) //ip+port :x.x.x.x:27015
	GetServerIPsBySteamID(serverSteamIds []string) (resp *GetServerIPsBySteamIDResp, err error)
}

type GameServersServiceConfig struct {
	ApiKey  string
	Timeout time.Duration
}
type gameServersServiceEntity struct {
	*GameServersServiceConfig
}

func NewGameServersService(cfg *GameServersServiceConfig) GameServersService {
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultReqTimeout
	}
	return &gameServersServiceEntity{GameServersServiceConfig: cfg}
}

func (s *gameServersServiceEntity) GetAccountList() (resp *AccountListResp, err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey}
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method:      http.MethodGet,
		Url:         UrlGetAccountList,
		Params:      params,
		ContentType: httpx.ContentJson,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	err = json.Unmarshal(httpResp.Content, &resp)
	return
}

func (s *gameServersServiceEntity) CreateAccount(appId int, memo string) (resp *CreatAccountResp, err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey}
	data := httpx.FormatFormData(map[string]string{"appid": util.IntToString(appId), "memo": memo})
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method:      http.MethodPost,
		Url:         UrlCreateAccount,
		Params:      params,
		Body:        data,
		ContentType: httpx.ContentForm,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	err = json.Unmarshal(httpResp.Content, &resp)
	return
}

func (s *gameServersServiceEntity) SetMemo(steamId string, memo string) (err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey}
	data := httpx.FormatFormData(map[string]string{"steamid": steamId, "memo": memo})
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method:      http.MethodPost,
		Url:         UrlSetMemo,
		Params:      params,
		Body:        data,
		ContentType: httpx.ContentForm,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	return
}

func (s *gameServersServiceEntity) ResetLoginToken(steamId string) (resp *ResetLoginTokenResp, err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey}
	data := httpx.FormatFormData(map[string]string{"steamid": steamId})
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method:      http.MethodPost,
		Url:         UrlResetLoginToken,
		Params:      params,
		Body:        data,
		ContentType: httpx.ContentForm,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	err = json.Unmarshal(httpResp.Content, &resp)
	return
}

func (s *gameServersServiceEntity) DeleteAccount(steamId string) (err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey}
	data := httpx.FormatFormData(map[string]string{"steamid": steamId})
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method:      http.MethodPost,
		Url:         UrlDeleteAccount,
		Params:      params,
		Body:        data,
		ContentType: httpx.ContentForm,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	return
}

func (s *gameServersServiceEntity) GetAccountPublicInfo(steamId string) (resp *GetAccountPublicInfoResp, err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey, "steamid": steamId}
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method: http.MethodGet,
		Url:    UrlGetAccountPublicInfo,
		Params: params,
		//Body:        data,
		ContentType: httpx.ContentForm,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	err = json.Unmarshal(httpResp.Content, &resp)
	return
}

func (s *gameServersServiceEntity) QueryLoginToken(loginToken string) (resp *QueryLoginTokenResp, err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey, "login_token": loginToken}
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method: http.MethodGet,
		Url:    UrlQueryLoginToken,
		Params: params,
		//Body:        data,
		ContentType: httpx.ContentForm,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	err = json.Unmarshal(httpResp.Content, &resp)
	return
}

// GetServerSteamIDsByIP ip:port
func (s *gameServersServiceEntity) GetServerSteamIDsByIP(serverIps []string) (resp *GetServerSteamIDsByIPResp, err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey}
	for index, sid := range serverIps {
		params[fmt.Sprintf("server_ips[%d]", index)] = sid
	}
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method:      http.MethodGet,
		Url:         UrlGetServerSteamIDsByIP,
		Params:      params,
		ContentType: httpx.ContentForm,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	err = json.Unmarshal(httpResp.Content, &resp)
	return
}

func (s *gameServersServiceEntity) GetServerIPsBySteamID(serverSteamIds []string) (resp *GetServerIPsBySteamIDResp, err error) {
	client := httpx.New(&httpx.Config{Timeout: s.Timeout})
	params := map[string]string{"key": s.ApiKey}
	for index, sid := range serverSteamIds {
		params[fmt.Sprintf("server_steamids[%d]", index)] = sid
	}
	httpResp, err := client.Send(context.Background(), &httpx.Request{
		Method:      http.MethodGet,
		Url:         UrlGetServerIPsBySteamID,
		Params:      params,
		ContentType: httpx.ContentForm,
	})
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("状态码异常，status: %d", httpResp.StatusCode)
		return
	}
	err = json.Unmarshal(httpResp.Content, &resp)
	return
}
