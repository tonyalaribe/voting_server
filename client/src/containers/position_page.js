import m from 'mithril';
import {PositionModel} from '../models/index.js';
import Nav from '../components/nav.js'

var PositionPage = {
  Position:"",
  Key:"",
  oncreate:function(vnode){
    console.log(vnode)
    console.log("vnode")
    console.log(vnode.attrs)
    let key = vnode.attrs.id
    console.log(key)
    this.Key = key
    PositionModel.GetPosition(key).then((resp)=>{
      let currentPosition = resp.Election.Positions.find((p)=>{
        console.log(p)
        console.log(key)
        return p.Key == key
      })
      console.log(currentPosition)
      this.Position = currentPosition.Title
      m.redraw()
    })


  },
  addContestant:function(e){
    e.preventDefault()
    let name = document.getElementById("contestantInput").value
    console.log(this.Key)
    PositionModel.NewCandidate(this.Key,name,this.Position)
  },
  view:function(){
    let contestants = PositionModel.Data.Contestants.map((c)=>{

      return m("p", m.trust(`>&nbsp;&nbsp;${c.Name}`))

    })
    return m("section.dt.w-100.vh-100.tc",
      m(Nav),
       m("div.dtc.v-mid",
        m("div.tl.w-50.bg-white.pa4.dib.shadow-4",
          m("a[href=/admin].link.pa2.black.db.bg-near-white",{oncreate:m.route.link},"<  go back"),

          m("div",
            m("h2",`Position: ${this.Position}`),
            m("div"),
            m("div",
              m("h3.fw4",
                m("span.dib.bb.pa1","Candidates")
              ),
              contestants
            ),
            m("form.tc",
              {
                onsubmit:this.addContestant.bind(this ),
              },
              m("input[type=text][placeholder=add contestant].pa2",{id:"contestantInput"}),
              m("button.pa2",{type:"submit"},"submit")
            )
          )
        )
      )
    )
  }
}

export default PositionPage;
