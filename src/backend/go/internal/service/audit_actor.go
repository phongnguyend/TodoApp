package service

func firstActor(actorUserID []*uint) *uint {
	if len(actorUserID) == 0 {
		return nil
	}
	return actorUserID[0]
}
