package network

// miscellaneous
const (
	HTTP_URL = "http://%s:%s%s"
)

// HTTP header
const (
	HDR_BDCAST = "Broadcast"
	HDR_ADDR   = "Peer_addr"
	HDR_PORT   = "Peer_port"
	HDR_HEIGHT = "Chain_height"
)

// main endpoints
const (
	NODE  = "/nodes"
	EMAIL = "/email"
	CHAIN = "/chain"
)

// email endpoints
const (
	EMAIL_NEW     = EMAIL + "/new"
	EMAIL_MAILBOX = EMAIL + "/mailbox"
	EMAIL_SENT    = EMAIL + "/sent"
)

//node endpoints
const (
	NODE_CONNECT      = NODE + "/connect"
	NODE_PENDINGMAILS = NODE + "/pending_mails"
	NODE_DISCONNECT   = NODE + "/disconnect"
	NODE_BROADCAST    = NODE + "/broadcast"
	NODE_PEERS        = NODE + "/peers"
	NODE_PEERS_UPDATE = NODE + "/peers/update"
)

// chain endpoints
const (
	CHAIN_GENERATE      = CHAIN + "/generate"
	CHAIN_ADD           = CHAIN + "/add"
	CHAIN_VIEW          = CHAIN + "/view"
	CHAIN_SYNC          = CHAIN + "/sync"
	CHAIN_GETBLOCK      = CHAIN + "/get_block"
	CHAIN_INFO          = CHAIN + "/info"
	CHAIN_MISSINGBLOCKS = CHAIN + "/missing_blocks"
)
