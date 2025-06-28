package state

import (
	privHandlers "server/state/privHandlers"
	pubHandlers "server/state/pubHandlers"
)

func (s *State) registerPublicRoutes() {
	// /public/* group
	{
		public := s.Router.Group("/public")

		// /public/ping
		public.GET("/ping", pubHandlers.Ping)
	}
}

func (s *State) registerPrivateRoutes() {
	// /admin/* group
	{
		admin := s.Router.Group("/admin")
		// /admin/ping group
		admin.GET("/ping", privHandlers.AdminPing)
		// /admin/pong group
		admin.GET("/pong", privHandlers.AdminPong)
	}

	// /ping
	s.Router.GET("/ping", privHandlers.Ping)
}
