package locations

func All() data.View {
	return data.Select(data.Type("location")).Sort(data.Get("name"))
}
