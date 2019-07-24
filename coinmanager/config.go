package coinmanager

type Config struct {
	InitialBalance    int
	UsedInvite        int
	InviteAccepted    int
	LikeReceived      int
	ReplyReceived     int
	FirstPostActivity int
	PopularPost       int
	LeaderboardFirst  int
	LeaderboardSecond int
	LeaderboardThird  int
	LeaderboardTopTen int
	LeaderboardRanked int
	BoughtThreadAlias int
	BoughtPostAlias   int
	BoughtChannel     int
}

func NewConfig() Config {
	return Config{}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	return nil
}
