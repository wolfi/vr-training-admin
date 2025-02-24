package pages

import (
	"github.com/a-h/templ"
	"github.com/saladinomario/vr-training-admin/templates/components"
	"github.com/saladinomario/vr-training-admin/templates/components/observers"
)

// ObserversIndex renders the observer management page
func ObserversIndex(observerList []observers.Observer) templ.Component {
	return templ.ComponentFunc(func(ctx templ.Context, w templ.Writer) error {
		if err := components.Layout("Observer Setup").Render(ctx, w); err != nil {
			return err
		}
		return observers.ObserverList(observerList).Render(ctx, w)
	})
}

// ObserverEdit renders the edit observer page
func ObserverEdit(observer observers.Observer) templ.Component {
	return templ.ComponentFunc(func(ctx templ.Context, w templ.Writer) error {
		if err := components.Layout("Edit Observer").Render(ctx, w); err != nil {
			return err
		}
		return observers.ObserverForm(observer).Render(ctx, w)
	})
}
