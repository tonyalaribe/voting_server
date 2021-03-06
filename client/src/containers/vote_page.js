import m from "mithril";
import Nav from '../components/nav.js';
import VoteModel from '../models/vote.js';


var VotePage = {}
VotePage.oncreate = function(){
  VoteModel.HasNotVoted()
  VoteModel.GetElectionData()
}

VotePage.CastVote = function(){
  let results = {}
  VoteModel.Positions.map((v)=>{
    results[v.Key] = document.getElementById(`position_${v.Key}`).value
  })
  VoteModel.SubmitVote(results)
}


VotePage.view = function(){
    return m("section.vh-100.w-100.tc",
      m(Nav),
      m("div.pt5",
      m("div.tc",
        m("img.h3.h4-ns",{src:"/public/images/logo.jpg"})
      ),
      m("h3.f5.navy", "College of Education Ikere-Ekiti"),
        m("h2.fw1","Cast your Vote")
      ),
      VoteModel.Positions.map((p)=>{
        return m("div.mw7.dib.tl.w-100",
          m("div.tl.bg-white.shadow-4.pa4.mv1",
            m("h2.fw1",`Position: ${p.Title}`),
            m("div",
              m("select.ba.bg-transparent",{id:`position_${p.Key}`},
                p.Contestants.map((contestant)=>{
                    return m("option",{value:contestant.ID},contestant.Name)
                })
              )
            )
          )
        )
      }),
        m("div.tc.mv5",
          m("button.pa3.bg-navy.white-80.shadow-4.bw0",{onclick:VotePage.CastVote},"Cast vote")
        )
      )

  }

export default VotePage;
