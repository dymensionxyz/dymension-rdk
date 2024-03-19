package types

func NewHub(hubId, channelId string) Hub {
	return Hub{
		HubId:     hubId,
		ChannelId: channelId,
	}
}
