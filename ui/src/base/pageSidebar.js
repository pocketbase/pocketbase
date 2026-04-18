window.app = window.app || {};
window.app.components = window.app.components || {};

const responsiveThreshold = 1000;

/**
 * A generic page sidebar component with resizable edge.
 *
 * @example
 * ```js
 * app.components.pageSidebar({},
 * t.nav({ className: "sidebar-content scrollable" },
 *     t.details({ className: "nav-group"},
 *         t.summary({ tabIndex: -1, onfocusout: () => false, onclick: () => false, onkeyup: () => false },
 *             "Group 1",
 *         ),
 *         t.a({ className: "nav-link", href: "..." }, "Link 1.1"),
 *         t.a({ className: "nav-link", href: "..." }, "Link 1.2"),
 *     ),
 *     t.details({ className: "nav-group"},
 *         t.summary({ tabIndex: -1, onfocusout: () => false, onclick: () => false, onkeyup: () => false },
 *             "Group 2",
 *         ),
 *         t.a({ className: "nav-link", href: "..." }, "Link 2.1"),
 *         t.a({ className: "nav-link", href: "..." }, "Link 2.2"),
 *     ),
 * )
 * ```
 *
 * @param  {Object} [propsArg]
 * @param  {Array<Element|Function>} [children]
 * @return {Element}
 */
window.app.components.pageSidebar = function(propsArg = {}, ...children) {
    let sidebarElem;

    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        widthHistoryKey: "pbPageSidebarWidth",
        onmount: undefined,
        onunmount: undefined,
    });

    const watchers = app.utils.extendStore(props, propsArg);

    const data = store({
        responsiveShow: false,
    });

    let responsiveBtn;
    function responsivePageSidebar() {
        if (!sidebarElem) {
            return;
        }

        if (window.innerWidth > responsiveThreshold) {
            data.responsiveShow = false;
            sidebarElem.dataset.responsive = false;
            responsiveBtn?.remove();
            responsiveBtn = null;
            return;
        }

        sidebarElem.dataset.responsive = true;

        if (!responsiveBtn) {
            responsiveBtn = t.button(
                {
                    type: "button",
                    className: "btn transparent secondary responsive-sidebar-btn",
                    title: "Toggle sidebar",
                    onclick: (e) => {
                        e.stopPropagation();
                        data.responsiveShow = !data.responsiveShow;
                    },
                },
                t.i({ className: "ri-menu-2-line", ariaHidden: true }),
            );

            document.body.querySelector(".page-header .breadcrumbs").before(responsiveBtn);
        }
    }

    function onOutsideClick(e) {
        if (e.target.closest(".responsive-close")) {
            data.responsiveShow = false;
            return;
        }

        if (
            e.target.closest(".page-sidebar")
            || e.target.closest(".app-header")
            || e.target.closest(".modal")
        ) {
            return; // inside click -> do nothing
        }

        e.preventDefault();
        e.stopImmediatePropagation();

        data.responsiveShow = false;

        return false;
    }

    watchers.push(
        watch(() => data.responsiveShow, (visible) => {
            if (visible) {
                window.addEventListener("click", onOutsideClick, true);
            } else {
                window.removeEventListener("click", onOutsideClick, true);
            }
        }),
    );

    sidebarElem = t.aside(
        {
            pbEvent: "pageSidebar",
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: () => `page-sidebar ${props.className} ${data.responsiveShow ? "active" : ""}`,
            onmount: (el) => {
                responsivePageSidebar(el);
                window.addEventListener("resize", responsivePageSidebar);
                props.onmount?.(el);
            },
            onunmount: (el) => {
                props.onunmount?.(el);

                window.removeEventListener("click", onOutsideClick, true);
                window.removeEventListener("resize", responsivePageSidebar);
                responsiveBtn?.remove();
                watchers.forEach((w) => w?.unwatch());
            },
        },
        (el) => {
            let sidebarWidth;

            if (props.widthHistoryKey) {
                sidebarWidth = localStorage.getItem(props.widthHistoryKey);
                if (sidebarWidth) {
                    el.style.width = sidebarWidth;
                }
            }

            return app.components.dragline({
                ondragstart: (e) => {
                    el._startWidth = el.offsetWidth;
                },
                ondragging: (e, diffX, diffY) => {
                    sidebarWidth = el._startWidth + diffX + "px";
                    el.style.width = sidebarWidth;

                    if (props.widthHistoryKey) {
                        localStorage.setItem(props.widthHistoryKey, sidebarWidth);
                    }
                },
            });
        },
        ...children,
    );

    return sidebarElem;
};
