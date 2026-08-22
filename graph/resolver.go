package graph

import "bugtracker/internal/services"

type Resolver struct {
	ProjectService *services.ProjectService
	BoardService   *services.BoardService
	CardService    *services.CardService
	MemberService  *services.MemberService
	UserService    *services.UserService
	MessageService *services.MessageService
}

type queryResolver struct {
	*Resolver
}

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
