import{S as Ce,i as Re,s as Ee,e as r,w as $,b as h,c as $e,f as m,g as f,h as n,m as we,x,aa as _e,ab as Pe,k as Te,ac as Be,n as Oe,t as ee,a as te,o as u,d as ge,ae as Ie,C as Ae,p as Me,r as z,u as Se,a9 as qe}from"./index-c45c880c.js";import{S as He}from"./SdkTabs-04dd5574.js";function ke(o,l,s){const a=o.slice();return a[6]=l[s],a}function he(o,l,s){const a=o.slice();return a[6]=l[s],a}function ve(o){let l;return{c(){l=r("p"),l.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",m(l,"class","txt-hint txt-sm txt-right")},m(s,a){f(s,l,a)},d(s){s&&u(l)}}}function ye(o,l){let s,a=l[6].code+"",v,i,c,p;function w(){return l[5](l[6])}return{key:o,first:null,c(){s=r("button"),v=$(a),i=h(),m(s,"class","tab-item"),z(s,"active",l[2]===l[6].code),this.first=s},m(b,g){f(b,s,g),n(s,v),n(s,i),c||(p=Se(s,"click",w),c=!0)},p(b,g){l=b,g&20&&z(s,"active",l[2]===l[6].code)},d(b){b&&u(s),c=!1,p()}}}function De(o,l){let s,a,v,i;return a=new qe({props:{content:l[6].body}}),{key:o,first:null,c(){s=r("div"),$e(a.$$.fragment),v=h(),m(s,"class","tab-item"),z(s,"active",l[2]===l[6].code),this.first=s},m(c,p){f(c,s,p),we(a,s,null),n(s,v),i=!0},p(c,p){l=c,(!i||p&20)&&z(s,"active",l[2]===l[6].code)},i(c){i||(ee(a.$$.fragment,c),i=!0)},o(c){te(a.$$.fragment,c),i=!1},d(c){c&&u(s),ge(a)}}}function Le(o){var ue,pe;let l,s,a=o[0].name+"",v,i,c,p,w,b,g,q=o[0].name+"",F,le,K,C,N,T,G,y,H,se,L,P,oe,J,U=o[0].name+"",Q,ae,V,ne,W,B,X,O,Y,I,Z,R,A,D=[],ie=new Map,ce,M,_=[],re=new Map,E;C=new He({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').delete('RECORD_ID');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').delete('RECORD_ID');
    `}});let k=o[1]&&ve(),j=o[4];const de=e=>e[6].code;for(let e=0;e<j.length;e+=1){let t=he(o,j,e),d=de(t);ie.set(d,D[e]=ye(d,t))}let S=o[4];const fe=e=>e[6].code;for(let e=0;e<S.length;e+=1){let t=ke(o,S,e),d=fe(t);re.set(d,_[e]=De(d,t))}return{c(){l=r("h3"),s=$("Delete ("),v=$(a),i=$(")"),c=h(),p=r("div"),w=r("p"),b=$("Delete a single "),g=r("strong"),F=$(q),le=$(" record."),K=h(),$e(C.$$.fragment),N=h(),T=r("h6"),T.textContent="API details",G=h(),y=r("div"),H=r("strong"),H.textContent="DELETE",se=h(),L=r("div"),P=r("p"),oe=$("/api/collections/"),J=r("strong"),Q=$(U),ae=$("/records/"),V=r("strong"),V.textContent=":id",ne=h(),k&&k.c(),W=h(),B=r("div"),B.textContent="Path parameters",X=h(),O=r("table"),O.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to delete.</td></tr></tbody>`,Y=h(),I=r("div"),I.textContent="Responses",Z=h(),R=r("div"),A=r("div");for(let e=0;e<D.length;e+=1)D[e].c();ce=h(),M=r("div");for(let e=0;e<_.length;e+=1)_[e].c();m(l,"class","m-b-sm"),m(p,"class","content txt-lg m-b-sm"),m(T,"class","m-b-xs"),m(H,"class","label label-primary"),m(L,"class","content"),m(y,"class","alert alert-danger"),m(B,"class","section-title"),m(O,"class","table-compact table-border m-b-base"),m(I,"class","section-title"),m(A,"class","tabs-header compact left"),m(M,"class","tabs-content"),m(R,"class","tabs")},m(e,t){f(e,l,t),n(l,s),n(l,v),n(l,i),f(e,c,t),f(e,p,t),n(p,w),n(w,b),n(w,g),n(g,F),n(w,le),f(e,K,t),we(C,e,t),f(e,N,t),f(e,T,t),f(e,G,t),f(e,y,t),n(y,H),n(y,se),n(y,L),n(L,P),n(P,oe),n(P,J),n(J,Q),n(P,ae),n(P,V),n(y,ne),k&&k.m(y,null),f(e,W,t),f(e,B,t),f(e,X,t),f(e,O,t),f(e,Y,t),f(e,I,t),f(e,Z,t),f(e,R,t),n(R,A);for(let d=0;d<D.length;d+=1)D[d]&&D[d].m(A,null);n(R,ce),n(R,M);for(let d=0;d<_.length;d+=1)_[d]&&_[d].m(M,null);E=!0},p(e,[t]){var me,be;(!E||t&1)&&a!==(a=e[0].name+"")&&x(v,a),(!E||t&1)&&q!==(q=e[0].name+"")&&x(F,q);const d={};t&9&&(d.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').delete('RECORD_ID');
    `),t&9&&(d.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').delete('RECORD_ID');
    `),C.$set(d),(!E||t&1)&&U!==(U=e[0].name+"")&&x(Q,U),e[1]?k||(k=ve(),k.c(),k.m(y,null)):k&&(k.d(1),k=null),t&20&&(j=e[4],D=_e(D,t,de,1,e,j,ie,A,Pe,ye,null,he)),t&20&&(S=e[4],Te(),_=_e(_,t,fe,1,e,S,re,M,Be,De,null,ke),Oe())},i(e){if(!E){ee(C.$$.fragment,e);for(let t=0;t<S.length;t+=1)ee(_[t]);E=!0}},o(e){te(C.$$.fragment,e);for(let t=0;t<_.length;t+=1)te(_[t]);E=!1},d(e){e&&u(l),e&&u(c),e&&u(p),e&&u(K),ge(C,e),e&&u(N),e&&u(T),e&&u(G),e&&u(y),k&&k.d(),e&&u(W),e&&u(B),e&&u(X),e&&u(O),e&&u(Y),e&&u(I),e&&u(Z),e&&u(R);for(let t=0;t<D.length;t+=1)D[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function Ue(o,l,s){let a,v,{collection:i=new Ie}=l,c=204,p=[];const w=b=>s(2,c=b.code);return o.$$set=b=>{"collection"in b&&s(0,i=b.collection)},o.$$.update=()=>{o.$$.dirty&1&&s(1,a=(i==null?void 0:i.deleteRule)===null),o.$$.dirty&3&&i!=null&&i.id&&(p.push({code:204,body:`
                null
            `}),p.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `}),a&&p.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}),p.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},s(3,v=Ae.getApiExampleUrl(Me.baseUrl)),[i,a,c,v,p,w]}class Fe extends Ce{constructor(l){super(),Re(this,l,Ue,Le,Ee,{collection:0})}}export{Fe as default};
