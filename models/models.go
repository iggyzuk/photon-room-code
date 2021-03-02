package models

// RoomCode is a server-side type that keeps track of all codes.
type RoomCode struct {
	Code    string
	Created bool
}

// CreateRoomRequest is a Photon type that is sent to us when a new room is created.
type CreateRoomRequest struct {
	ActorNr       int    `json:"ActorNr"`
	AppVersion    string `json:"AppVersion"`
	AppID         string `json:"AppId"`
	CreateOptions struct {
		MaxPlayers       int         `json:"MaxPlayers"`
		IsVisible        bool        `json:"IsVisible"`
		LobbyID          interface{} `json:"LobbyId"`
		LobbyType        int         `json:"LobbyType"`
		CustomProperties struct {
		} `json:"CustomProperties"`
		EmptyRoomTTL       int         `json:"EmptyRoomTTL"`
		PlayerTTL          int         `json:"PlayerTTL"`
		CheckUserOnJoin    bool        `json:"CheckUserOnJoin"`
		DeleteCacheOnLeave bool        `json:"DeleteCacheOnLeave"`
		SuppressRoomEvents bool        `json:"SuppressRoomEvents"`
		PublishUserID      bool        `json:"PublishUserId"`
		ExpectedUsers      interface{} `json:"ExpectedUsers"`
	} `json:"CreateOptions"`
	GameID   string `json:"GameId"`
	Region   string `json:"Region"`
	Type     string `json:"Type"`
	UserID   string `json:"UserId"`
	Nickname string `json:"Nickname"`
}

// CloseRoomRequest is a Photon type that is sent to us when a room is closed.
type CloseRoomRequest struct {
	ActorCount int    `json:"ActorCount"`
	AppVersion string `json:"AppVersion"`
	AppID      string `json:"AppId"`
	GameID     string `json:"GameId"`
	Region     string `json:"Region"`
	Type       string `json:"Type"`
}

// RoomResponse is what's send back to Photon.
type RoomResponse struct {
	State      string `json:"State"`
	ResultCode int    `json:"ResultCode"`
}
