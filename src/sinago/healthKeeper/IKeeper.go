package healthKeeper

type IKeeper interface {

	Start()

	SetChecker(interface{})

	OnStart(interface{})

	OnStatus(interface{})

	List(int)

}
