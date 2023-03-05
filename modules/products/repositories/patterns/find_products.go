package patterns

type IFindProductBuidler interface {
}

type findProductBuilder struct {
	query string
}

func FindProductBuilder() IFindProductBuidler {
	return &findProductBuilder{}
}

type findProductEngineer struct {
	builder IFindProductBuidler
}

func FindProductEngineer(b IFindProductBuidler) *findProductEngineer {
	return &findProductEngineer{builder: b}
}

func (en *findProductEngineer) FindProduct() IFindProductBuidler {
	return nil
}

func (en *findProductEngineer) FindOneProduct() IFindProductBuidler {
	return nil
}
