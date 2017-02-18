import m from "mithril"

var Nav = {
  logout:function(){
    localStorage.removeItem("auth")
  },
  view:function(){
    return m("nav.fixed.w-100.bg-black.white.pa2.tc",
      m("a[href=/].white.link.dib.pa2",{oncreate: m.route.link},"home"),
      m("a[href=/results].white.link.dib.pa2",{oncreate: m.route.link},"results"),
      m("a.pointer.white.link.dib.pa2",{onclick:this.logout},"logout")
    )
  }
}

export default Nav;
