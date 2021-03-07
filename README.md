# Unique Room Codes API with Photon

This is a web API that creates unique room codes for multiplayer Unity games that use [Photon PUN](https://www.photonengine.com/pun).

When you want to play with friends online, you will need to create a private room, specify a unique name and a password. I think that's a lot of steps for something that can be solved with an automatically generated short unique room code. 

You press create room, get a code, and share it with friends.

The API generates four digit codes from 0000 to 9999.

- 0055
- 0232
- 1520
- 7002
- 9206

# Setup Photon Webhooks

You must setup Photon webhooks for the API to function properly. We're interested in webhooks that fire when a new room is created and when a room is closed – PathCreate and PathClose.

1. Go to the [Photon dashboard](https://dashboard.photonengine.com/)
2. Find your project
3. Press on manage
4. Go down to webhooks
5. Press edit
6. Set the details as following

- AsyncJoin: true
- BaseURL: https://example.herokuapp.com
- HasErrorInfo: true
- IsPersistent: true
- PathClose: room/close
- PathCreate room/create

# Handlers

The API contains three handlers:

- /room/gen_code
- /room/create
- /room/close

## /room/gen_code

Generates and returns a unique room code. It must be called before creating a room because the code will be used as the name of the room.

### Code Example

```cs
// Ask the server for a new code.
var codeRequest = UnityWebRequest.Get($"{ServerURL}/room/gen_code/");
var codeData = await codeRequest.SendWebRequest();
var code = codeData.downloadHandler.text;

// Create a new room.
var roomOptions = new RoomOptions
{
	IsVisible = true,
	MaxPlayers = MaxPlayers
};

PhotonNetwork.CreateRoom(code, roomOptions, TypedLobby.Default);
```

## /room/create

Called by a Photon when a new room is created. We use it because we need to confirm that a room was created. Unconfirmed codes timeout after 10 seconds.

## /room/close

Called by a Photon when a room is closed. We use it to release a code, allow it to be used again by another player.

# Host
You can easily host it on [Heroku](https://heroku.com/).

# Future

Possibly need to separate the codes by region and app version. Feel free to contribute.
