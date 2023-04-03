import{S as Ze,i as et,s as tt,N as Ye,e as o,w as m,b as u,c as _e,f as _,g as r,h as l,m as ke,x as me,O as Ve,P as lt,k as st,Q as nt,n as ot,t as z,a as G,o as d,d as he,T as it,C as ze,p as at,r as J,u as rt}from"./index-4eea3e34.js";import{S as dt}from"./SdkTabs-5d6cc1d4.js";function Ge(i,s,n){const a=i.slice();return a[6]=s[n],a}function Je(i,s,n){const a=i.slice();return a[6]=s[n],a}function Ke(i){let s;return{c(){s=o("p"),s.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",_(s,"class","txt-hint txt-sm txt-right")},m(n,a){r(n,s,a)},d(n){n&&d(s)}}}function We(i,s){let n,a=s[6].code+"",y,c,f,b;function F(){return s[5](s[6])}return{key:i,first:null,c(){n=o("button"),y=m(a),c=u(),_(n,"class","tab-item"),J(n,"active",s[2]===s[6].code),this.first=n},m(h,R){r(h,n,R),l(n,y),l(n,c),f||(b=rt(n,"click",F),f=!0)},p(h,R){s=h,R&20&&J(n,"active",s[2]===s[6].code)},d(h){h&&d(n),f=!1,b()}}}function Xe(i,s){let n,a,y,c;return a=new Ye({props:{content:s[6].body}}),{key:i,first:null,c(){n=o("div"),_e(a.$$.fragment),y=u(),_(n,"class","tab-item"),J(n,"active",s[2]===s[6].code),this.first=n},m(f,b){r(f,n,b),ke(a,n,null),l(n,y),c=!0},p(f,b){s=f,(!c||b&20)&&J(n,"active",s[2]===s[6].code)},i(f){c||(z(a.$$.fragment,f),c=!0)},o(f){G(a.$$.fragment,f),c=!1},d(f){f&&d(n),he(a)}}}function ct(i){var Ne,Ue;let s,n,a=i[0].name+"",y,c,f,b,F,h,R,N=i[0].name+"",K,ve,W,g,X,B,Y,$,U,we,j,E,ye,Z,Q=i[0].name+"",ee,$e,te,Ce,le,x,se,A,ne,I,oe,O,ie,Fe,ae,T,re,Re,de,ge,k,Oe,S,Te,De,Pe,ce,Ee,fe,Se,Be,xe,pe,Ae,ue,M,be,D,H,C=[],Ie=new Map,Me,q,v=[],He=new Map,P;g=new dt({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        const record = await pb.collection('${(Ne=i[0])==null?void 0:Ne.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        final record = await pb.collection('${(Ue=i[0])==null?void 0:Ue.name}').getOne('RECORD_ID',
          'expand': 'relField1,relField2.subRelField',
        );
    `}});let w=i[1]&&Ke();S=new Ye({props:{content:"?expand=relField1,relField2.subRelField"}});let V=i[4];const qe=e=>e[6].code;for(let e=0;e<V.length;e+=1){let t=Je(i,V,e),p=qe(t);Ie.set(p,C[e]=We(p,t))}let L=i[4];const Le=e=>e[6].code;for(let e=0;e<L.length;e+=1){let t=Ge(i,L,e),p=Le(t);He.set(p,v[e]=Xe(p,t))}return{c(){s=o("h3"),n=m("View ("),y=m(a),c=m(")"),f=u(),b=o("div"),F=o("p"),h=m("Fetch a single "),R=o("strong"),K=m(N),ve=m(" record."),W=u(),_e(g.$$.fragment),X=u(),B=o("h6"),B.textContent="API details",Y=u(),$=o("div"),U=o("strong"),U.textContent="GET",we=u(),j=o("div"),E=o("p"),ye=m("/api/collections/"),Z=o("strong"),ee=m(Q),$e=m("/records/"),te=o("strong"),te.textContent=":id",Ce=u(),w&&w.c(),le=u(),x=o("div"),x.textContent="Path Parameters",se=u(),A=o("table"),A.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to view.</td></tr></tbody>`,ne=u(),I=o("div"),I.textContent="Query parameters",oe=u(),O=o("table"),ie=o("thead"),ie.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Fe=u(),ae=o("tbody"),T=o("tr"),re=o("td"),re.textContent="expand",Re=u(),de=o("td"),de.innerHTML='<span class="label">String</span>',ge=u(),k=o("td"),Oe=m(`Auto expand record relations. Ex.:
                `),_e(S.$$.fragment),Te=m(`
                Supports up to 6-levels depth nested relations expansion. `),De=o("br"),Pe=m(`
                The expanded relations will be appended to the record under the
                `),ce=o("code"),ce.textContent="expand",Ee=m(" property (eg. "),fe=o("code"),fe.textContent='"expand": {"relField1": {...}, ...}',Se=m(`).
                `),Be=o("br"),xe=m(`
                Only the relations to which the request user has permissions to `),pe=o("strong"),pe.textContent="view",Ae=m(" will be expanded."),ue=u(),M=o("div"),M.textContent="Responses",be=u(),D=o("div"),H=o("div");for(let e=0;e<C.length;e+=1)C[e].c();Me=u(),q=o("div");for(let e=0;e<v.length;e+=1)v[e].c();_(s,"class","m-b-sm"),_(b,"class","content txt-lg m-b-sm"),_(B,"class","m-b-xs"),_(U,"class","label label-primary"),_(j,"class","content"),_($,"class","alert alert-info"),_(x,"class","section-title"),_(A,"class","table-compact table-border m-b-base"),_(I,"class","section-title"),_(O,"class","table-compact table-border m-b-base"),_(M,"class","section-title"),_(H,"class","tabs-header compact left"),_(q,"class","tabs-content"),_(D,"class","tabs")},m(e,t){r(e,s,t),l(s,n),l(s,y),l(s,c),r(e,f,t),r(e,b,t),l(b,F),l(F,h),l(F,R),l(R,K),l(F,ve),r(e,W,t),ke(g,e,t),r(e,X,t),r(e,B,t),r(e,Y,t),r(e,$,t),l($,U),l($,we),l($,j),l(j,E),l(E,ye),l(E,Z),l(Z,ee),l(E,$e),l(E,te),l($,Ce),w&&w.m($,null),r(e,le,t),r(e,x,t),r(e,se,t),r(e,A,t),r(e,ne,t),r(e,I,t),r(e,oe,t),r(e,O,t),l(O,ie),l(O,Fe),l(O,ae),l(ae,T),l(T,re),l(T,Re),l(T,de),l(T,ge),l(T,k),l(k,Oe),ke(S,k,null),l(k,Te),l(k,De),l(k,Pe),l(k,ce),l(k,Ee),l(k,fe),l(k,Se),l(k,Be),l(k,xe),l(k,pe),l(k,Ae),r(e,ue,t),r(e,M,t),r(e,be,t),r(e,D,t),l(D,H);for(let p=0;p<C.length;p+=1)C[p]&&C[p].m(H,null);l(D,Me),l(D,q);for(let p=0;p<v.length;p+=1)v[p]&&v[p].m(q,null);P=!0},p(e,[t]){var je,Qe;(!P||t&1)&&a!==(a=e[0].name+"")&&me(y,a),(!P||t&1)&&N!==(N=e[0].name+"")&&me(K,N);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const record = await pb.collection('${(je=e[0])==null?void 0:je.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final record = await pb.collection('${(Qe=e[0])==null?void 0:Qe.name}').getOne('RECORD_ID',
          'expand': 'relField1,relField2.subRelField',
        );
    `),g.$set(p),(!P||t&1)&&Q!==(Q=e[0].name+"")&&me(ee,Q),e[1]?w||(w=Ke(),w.c(),w.m($,null)):w&&(w.d(1),w=null),t&20&&(V=e[4],C=Ve(C,t,qe,1,e,V,Ie,H,lt,We,null,Je)),t&20&&(L=e[4],st(),v=Ve(v,t,Le,1,e,L,He,q,nt,Xe,null,Ge),ot())},i(e){if(!P){z(g.$$.fragment,e),z(S.$$.fragment,e);for(let t=0;t<L.length;t+=1)z(v[t]);P=!0}},o(e){G(g.$$.fragment,e),G(S.$$.fragment,e);for(let t=0;t<v.length;t+=1)G(v[t]);P=!1},d(e){e&&d(s),e&&d(f),e&&d(b),e&&d(W),he(g,e),e&&d(X),e&&d(B),e&&d(Y),e&&d($),w&&w.d(),e&&d(le),e&&d(x),e&&d(se),e&&d(A),e&&d(ne),e&&d(I),e&&d(oe),e&&d(O),he(S),e&&d(ue),e&&d(M),e&&d(be),e&&d(D);for(let t=0;t<C.length;t+=1)C[t].d();for(let t=0;t<v.length;t+=1)v[t].d()}}}function ft(i,s,n){let a,y,{collection:c=new it}=s,f=200,b=[];const F=h=>n(2,f=h.code);return i.$$set=h=>{"collection"in h&&n(0,c=h.collection)},i.$$.update=()=>{i.$$.dirty&1&&n(1,a=(c==null?void 0:c.viewRule)===null),i.$$.dirty&3&&c!=null&&c.id&&(b.push({code:200,body:JSON.stringify(ze.dummyCollectionRecord(c),null,2)}),a&&b.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}),b.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},n(3,y=ze.getApiExampleUrl(at.baseUrl)),[c,a,f,y,b,F]}class bt extends Ze{constructor(s){super(),et(this,s,ft,ct,tt,{collection:0})}}export{bt as default};
