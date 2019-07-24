package globalid

type Fixture struct {
	ids   []ID
	index int
}

func NewFixture() *Fixture {
	return &Fixture{ids: []ID{
		"7906d4df-f636-4c7a-8503-0b6212f6d508",
		"435c4647-8a2e-45f5-9e70-247743fbb974",
		"b6b4d35a-a07d-48a5-8044-7c12f9a8e5d3",
		"441ed35d-64a4-409f-9626-599cf0285412",
		"e6e73410-a668-41bb-8ee8-1c7c616ceee6",
		"5e033154-4304-42d3-aef7-b847122d3080",
		"88b245d9-5593-432f-bccb-0224eb4d93f1",
		"1695d294-8267-4ab4-be4c-c2500bad8363",
		"6fc7bab0-58cc-47ef-b171-f84aa7ee8024",
		"7bd8b9ca-a3fd-4a0d-8d9c-18a9fc37284a",
		"823ed211-d54c-48fb-a470-613e2903286a",
		"3b9b3d1f-14fe-4674-a278-60ee5d6ae779",
		"26977874-05d0-4b23-819a-9c8aa8e14220",
		"3bb20aa6-dc52-4671-858b-0db9610aad92",
		"5fb41610-296c-4fb7-8b9f-075fe50015e4",
		"7a7e1460-bb8a-4449-aad2-2eb63584a6d6",
		"fc5e8d04-e82e-498f-9197-904a8309e01b",
		"967ae62c-608d-498c-96c3-b2056e35998d",
		"dfc745f8-1653-4072-8763-40e3fbbfb286",
		"d780f484-e827-4669-82dd-7fb9857b183f",
		"370c1a3c-6736-412a-a41d-fc548c398b7f",
		"8472a872-72ca-421c-a2b2-4efedc162a86",
		"b707305f-14f6-4bd8-89da-744470f09a74",
		"5dcb2f1c-c821-447b-850c-6478f707d9ed",
		"2b832de1-5f5a-4411-9b39-78b8afd942e2",
		"32633b5c-af43-4798-8a15-73cc36870513",
		"3f53ee37-b6e0-4f25-8d0a-bb4f92c4d941",
		"a0bce959-a3a8-476f-bc4e-af94cf83cfca",
		"2845c0af-c6f6-4ec7-ba4a-e5817aeda420",
		"bafb1a77-853c-47ea-b547-4a78fc00bab4",
		"837b4179-67ce-487c-8517-51ea6e3ccce3",
		"db0a2f0e-08d4-49f4-a92f-bbe55ea1e1c8",
		"52006ce8-e177-4c76-87ef-9c89d60d16e1",
		"11f3f8f9-d4cd-4ba0-8643-c3351f741dda",
		"a26e6f32-7e7b-4ac8-b009-e97400c1d753",
		"c51d0fae-9995-4497-a039-7c3ab9b11bdd",
		"926adc13-c46e-4540-86c4-2b9fd85debbb",
		"387e6bce-c7f0-437c-b83a-4a42543f1b81",
	}}
}

func (f *Fixture) Next() ID {
	defer func() { f.index++ }()
	return f.ids[f.index]
}

func (f *Fixture) Fixed(i int) ID {
	return f.ids[i]
}
