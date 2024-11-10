package engine

func (g *Game) Update() {}

func (s *Scene) AddComponent(component UIComponent) {
	s.Components = append(s.Components, component)
}
