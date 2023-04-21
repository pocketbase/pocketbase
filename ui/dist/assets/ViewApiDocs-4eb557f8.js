import{S as Ze,i as et,s as tt,M as Ye,e as o,w as m,b as u,c as _e,f as _,g as r,h as l,m as ke,x as me,N as ze,O as lt,k as st,P as nt,n as ot,t as G,a as J,o as d,d as he,T as it,C as Ge,p as at,r as K,u as rt}from"./index-43742f47.js";import{S as dt}from"./SdkTabs-292a1bcb.js";function Je(i,s,n){const a=i.slice();return a[6]=s[n],a}function Ke(i,s,n){const a=i.slice();return a[6]=s[n],a}function Qe(i){let s;return{c(){s=o("p"),s.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",_(s,"class","txt-hint txt-sm txt-right")},m(n,a){r(n,s,a)},d(n){n&&d(s)}}}function We(i,s){let n,a=s[6].code+"",y,c,f,b;function F(){return s[5](s[6])}return{key:i,first:null,c(){n=o("button"),y=m(a),c=u(),_(n,"class","tab-item"),K(n,"active",s[2]===s[6].code),this.first=n},m(h,R){r(h,n,R),l(n,y),l(n,c),f||(b=rt(n,"click",F),f=!0)},p(h,R){s=h,R&20&&K(n,"active",s[2]===s[6].code)},d(h){h&&d(n),f=!1,b()}}}function Xe(i,s){let n,a,y,c;return a=new Ye({props:{content:s[6].body}}),{key:i,first:null,c(){n=o("div"),_e(a.$$.fragment),y=u(),_(n,"class","tab-item"),K(n,"active",s[2]===s[6].code),this.first=n},m(f,b){r(f,n,b),ke(a,n,null),l(n,y),c=!0},p(f,b){s=f,(!c||b&20)&&K(n,"active",s[2]===s[6].code)},i(f){c||(G(a.$$.fragment,f),c=!0)},o(f){J(a.$$.fragment,f),c=!1},d(f){f&&d(n),he(a)}}}function ct(i){var Ne,Ue;let s,n,a=i[0].name+"",y,c,f,b,F,h,R,N=i[0].name+"",Q,ve,W,g,X,B,Y,$,U,we,j,E,ye,Z,V=i[0].name+"",ee,$e,te,Ce,le,M,se,x,ne,A,oe,O,ie,Fe,ae,T,re,Re,de,ge,k,Oe,S,Te,De,Pe,ce,Ee,fe,Se,Be,Me,pe,xe,ue,I,be,D,H,C=[],Ae=new Map,Ie,q,v=[],He=new Map,P;g=new dt({props:{js:`
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
          expand: 'relField1,relField2.subRelField',
        );
    `}});let w=i[1]&&Qe();S=new Ye({props:{content:"?expand=relField1,relField2.subRelField"}});let z=i[4];const qe=e=>e[6].code;for(let e=0;e<z.length;e+=1){let t=Ke(i,z,e),p=qe(t);Ae.set(p,C[e]=We(p,t))}let L=i[4];const Le=e=>e[6].code;for(let e=0;e<L.length;e+=1){let t=Je(i,L,e),p=Le(t);He.set(p,v[e]=Xe(p,t))}return{c(){s=o("h3"),n=m("View ("),y=m(a),c=m(")"),f=u(),b=o("div"),F=o("p"),h=m("Fetch a single "),R=o("strong"),Q=m(N),ve=m(" record."),W=u(),_e(g.$$.fragment),X=u(),B=o("h6"),B.textContent="API details",Y=u(),$=o("div"),U=o("strong"),U.textContent="GET",we=u(),j=o("div"),E=o("p"),ye=m("/api/collections/"),Z=o("strong"),ee=m(V),$e=m("/records/"),te=o("strong"),te.textContent=":id",Ce=u(),w&&w.c(),le=u(),M=o("div"),M.textContent="Path Parameters",se=u(),x=o("table"),x.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to view.</td></tr></tbody>`,ne=u(),A=o("div"),A.textContent="Query parameters",oe=u(),O=o("table"),ie=o("thead"),ie.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Fe=u(),ae=o("tbody"),T=o("tr"),re=o("td"),re.textContent="expand",Re=u(),de=o("td"),de.innerHTML='<span class="label">String</span>',ge=u(),k=o("td"),Oe=m(`Auto expand record relations. Ex.:
                `),_e(S.$$.fragment),Te=m(`
                Supports up to 6-levels depth nested relations expansion. `),De=o("br"),Pe=m(`
                The expanded relations will be appended to the record under the
                `),ce=o("code"),ce.textContent="expand",Ee=m(" property (eg. "),fe=o("code"),fe.textContent='"expand": {"relField1": {...}, ...}',Se=m(`).
                `),Be=o("br"),Me=m(`
                Only the relations to which the request user has permissions to `),pe=o("strong"),pe.textContent="view",xe=m(" will be expanded."),ue=u(),I=o("div"),I.textContent="Responses",be=u(),D=o("div"),H=o("div");for(let e=0;e<C.length;e+=1)C[e].c();Ie=u(),q=o("div");for(let e=0;e<v.length;e+=1)v[e].c();_(s,"class","m-b-sm"),_(b,"class","content txt-lg m-b-sm"),_(B,"class","m-b-xs"),_(U,"class","label label-primary"),_(j,"class","content"),_($,"class","alert alert-info"),_(M,"class","section-title"),_(x,"class","table-compact table-border m-b-base"),_(A,"class","section-title"),_(O,"class","table-compact table-border m-b-base"),_(I,"class","section-title"),_(H,"class","tabs-header compact left"),_(q,"class","tabs-content"),_(D,"class","tabs")},m(e,t){r(e,s,t),l(s,n),l(s,y),l(s,c),r(e,f,t),r(e,b,t),l(b,F),l(F,h),l(F,R),l(R,Q),l(F,ve),r(e,W,t),ke(g,e,t),r(e,X,t),r(e,B,t),r(e,Y,t),r(e,$,t),l($,U),l($,we),l($,j),l(j,E),l(E,ye),l(E,Z),l(Z,ee),l(E,$e),l(E,te),l($,Ce),w&&w.m($,null),r(e,le,t),r(e,M,t),r(e,se,t),r(e,x,t),r(e,ne,t),r(e,A,t),r(e,oe,t),r(e,O,t),l(O,ie),l(O,Fe),l(O,ae),l(ae,T),l(T,re),l(T,Re),l(T,de),l(T,ge),l(T,k),l(k,Oe),ke(S,k,null),l(k,Te),l(k,De),l(k,Pe),l(k,ce),l(k,Ee),l(k,fe),l(k,Se),l(k,Be),l(k,Me),l(k,pe),l(k,xe),r(e,ue,t),r(e,I,t),r(e,be,t),r(e,D,t),l(D,H);for(let p=0;p<C.length;p+=1)C[p]&&C[p].m(H,null);l(D,Ie),l(D,q);for(let p=0;p<v.length;p+=1)v[p]&&v[p].m(q,null);P=!0},p(e,[t]){var je,Ve;(!P||t&1)&&a!==(a=e[0].name+"")&&me(y,a),(!P||t&1)&&N!==(N=e[0].name+"")&&me(Q,N);const p={};t&9&&(p.js=`
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

        final record = await pb.collection('${(Ve=e[0])==null?void 0:Ve.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `),g.$set(p),(!P||t&1)&&V!==(V=e[0].name+"")&&me(ee,V),e[1]?w||(w=Qe(),w.c(),w.m($,null)):w&&(w.d(1),w=null),t&20&&(z=e[4],C=ze(C,t,qe,1,e,z,Ae,H,lt,We,null,Ke)),t&20&&(L=e[4],st(),v=ze(v,t,Le,1,e,L,He,q,nt,Xe,null,Je),ot())},i(e){if(!P){G(g.$$.fragment,e),G(S.$$.fragment,e);for(let t=0;t<L.length;t+=1)G(v[t]);P=!0}},o(e){J(g.$$.fragment,e),J(S.$$.fragment,e);for(let t=0;t<v.length;t+=1)J(v[t]);P=!1},d(e){e&&d(s),e&&d(f),e&&d(b),e&&d(W),he(g,e),e&&d(X),e&&d(B),e&&d(Y),e&&d($),w&&w.d(),e&&d(le),e&&d(M),e&&d(se),e&&d(x),e&&d(ne),e&&d(A),e&&d(oe),e&&d(O),he(S),e&&d(ue),e&&d(I),e&&d(be),e&&d(D);for(let t=0;t<C.length;t+=1)C[t].d();for(let t=0;t<v.length;t+=1)v[t].d()}}}function ft(i,s,n){let a,y,{collection:c=new it}=s,f=200,b=[];const F=h=>n(2,f=h.code);return i.$$set=h=>{"collection"in h&&n(0,c=h.collection)},i.$$.update=()=>{i.$$.dirty&1&&n(1,a=(c==null?void 0:c.viewRule)===null),i.$$.dirty&3&&c!=null&&c.id&&(b.push({code:200,body:JSON.stringify(Ge.dummyCollectionRecord(c),null,2)}),a&&b.push({code:403,body:`
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
            `}))},n(3,y=Ge.getApiExampleUrl(at.baseUrl)),[c,a,f,y,b,F]}class bt extends Ze{constructor(s){super(),et(this,s,ft,ct,tt,{collection:0})}}export{bt as default};
