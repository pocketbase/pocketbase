import{S as lt,i as st,s as nt,V as at,W as tt,X as K,h as r,d as W,t as V,a as j,I as ve,Z as Ge,_ as ot,C as it,$ as rt,D as dt,l as d,n as l,m as X,u as a,A as _,v as b,c as Z,w as m,J as Ke,p as ct,k as Y,o as pt}from"./index-CkK5VYgS.js";import{F as ut}from"./FieldsQueryParam-Z-S0qGe1.js";function We(o,s,n){const i=o.slice();return i[6]=s[n],i}function Xe(o,s,n){const i=o.slice();return i[6]=s[n],i}function Ze(o){let s;return{c(){s=a("p"),s.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",m(s,"class","txt-hint txt-sm txt-right")},m(n,i){d(n,s,i)},d(n){n&&r(s)}}}function Ye(o,s){let n,i,v;function p(){return s[5](s[6])}return{key:o,first:null,c(){n=a("button"),n.textContent=`${s[6].code} `,m(n,"class","tab-item"),Y(n,"active",s[2]===s[6].code),this.first=n},m(c,f){d(c,n,f),i||(v=pt(n,"click",p),i=!0)},p(c,f){s=c,f&20&&Y(n,"active",s[2]===s[6].code)},d(c){c&&r(n),i=!1,v()}}}function et(o,s){let n,i,v,p;return i=new tt({props:{content:s[6].body}}),{key:o,first:null,c(){n=a("div"),Z(i.$$.fragment),v=b(),m(n,"class","tab-item"),Y(n,"active",s[2]===s[6].code),this.first=n},m(c,f){d(c,n,f),X(i,n,null),l(n,v),p=!0},p(c,f){s=c,(!p||f&20)&&Y(n,"active",s[2]===s[6].code)},i(c){p||(j(i.$$.fragment,c),p=!0)},o(c){V(i.$$.fragment,c),p=!1},d(c){c&&r(n),W(i)}}}function ft(o){var Je,Ne;let s,n,i=o[0].name+"",v,p,c,f,w,C,ee,J=o[0].name+"",te,$e,le,F,se,B,ne,$,N,ye,Q,T,we,ae,z=o[0].name+"",oe,Ce,ie,Fe,re,I,de,S,ce,x,pe,R,ue,Re,M,D,fe,De,be,Oe,h,Pe,E,Te,Ee,Ae,me,Be,_e,Ie,Se,xe,he,Me,qe,A,ke,q,ge,O,H,y=[],He=new Map,Le,L,k=[],Ue=new Map,P;F=new at({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        const record = await pb.collection('${(Je=o[0])==null?void 0:Je.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        final record = await pb.collection('${(Ne=o[0])==null?void 0:Ne.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let g=o[1]&&Ze();E=new tt({props:{content:"?expand=relField1,relField2.subRelField"}}),A=new ut({});let G=K(o[4]);const Ve=e=>e[6].code;for(let e=0;e<G.length;e+=1){let t=Xe(o,G,e),u=Ve(t);He.set(u,y[e]=Ye(u,t))}let U=K(o[4]);const je=e=>e[6].code;for(let e=0;e<U.length;e+=1){let t=We(o,U,e),u=je(t);Ue.set(u,k[e]=et(u,t))}return{c(){s=a("h3"),n=_("View ("),v=_(i),p=_(")"),c=b(),f=a("div"),w=a("p"),C=_("Fetch a single "),ee=a("strong"),te=_(J),$e=_(" record."),le=b(),Z(F.$$.fragment),se=b(),B=a("h6"),B.textContent="API details",ne=b(),$=a("div"),N=a("strong"),N.textContent="GET",ye=b(),Q=a("div"),T=a("p"),we=_("/api/collections/"),ae=a("strong"),oe=_(z),Ce=_("/records/"),ie=a("strong"),ie.textContent=":id",Fe=b(),g&&g.c(),re=b(),I=a("div"),I.textContent="Path Parameters",de=b(),S=a("table"),S.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to view.</td></tr></tbody>',ce=b(),x=a("div"),x.textContent="Query parameters",pe=b(),R=a("table"),ue=a("thead"),ue.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Re=b(),M=a("tbody"),D=a("tr"),fe=a("td"),fe.textContent="expand",De=b(),be=a("td"),be.innerHTML='<span class="label">String</span>',Oe=b(),h=a("td"),Pe=_(`Auto expand record relations. Ex.:
                `),Z(E.$$.fragment),Te=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ee=a("br"),Ae=_(`
                The expanded relations will be appended to the record under the
                `),me=a("code"),me.textContent="expand",Be=_(" property (eg. "),_e=a("code"),_e.textContent='"expand": {"relField1": {...}, ...}',Ie=_(`).
                `),Se=a("br"),xe=_(`
                Only the relations to which the request user has permissions to `),he=a("strong"),he.textContent="view",Me=_(" will be expanded."),qe=b(),Z(A.$$.fragment),ke=b(),q=a("div"),q.textContent="Responses",ge=b(),O=a("div"),H=a("div");for(let e=0;e<y.length;e+=1)y[e].c();Le=b(),L=a("div");for(let e=0;e<k.length;e+=1)k[e].c();m(s,"class","m-b-sm"),m(f,"class","content txt-lg m-b-sm"),m(B,"class","m-b-xs"),m(N,"class","label label-primary"),m(Q,"class","content"),m($,"class","alert alert-info"),m(I,"class","section-title"),m(S,"class","table-compact table-border m-b-base"),m(x,"class","section-title"),m(R,"class","table-compact table-border m-b-base"),m(q,"class","section-title"),m(H,"class","tabs-header compact combined left"),m(L,"class","tabs-content"),m(O,"class","tabs")},m(e,t){d(e,s,t),l(s,n),l(s,v),l(s,p),d(e,c,t),d(e,f,t),l(f,w),l(w,C),l(w,ee),l(ee,te),l(w,$e),d(e,le,t),X(F,e,t),d(e,se,t),d(e,B,t),d(e,ne,t),d(e,$,t),l($,N),l($,ye),l($,Q),l(Q,T),l(T,we),l(T,ae),l(ae,oe),l(T,Ce),l(T,ie),l($,Fe),g&&g.m($,null),d(e,re,t),d(e,I,t),d(e,de,t),d(e,S,t),d(e,ce,t),d(e,x,t),d(e,pe,t),d(e,R,t),l(R,ue),l(R,Re),l(R,M),l(M,D),l(D,fe),l(D,De),l(D,be),l(D,Oe),l(D,h),l(h,Pe),X(E,h,null),l(h,Te),l(h,Ee),l(h,Ae),l(h,me),l(h,Be),l(h,_e),l(h,Ie),l(h,Se),l(h,xe),l(h,he),l(h,Me),l(M,qe),X(A,M,null),d(e,ke,t),d(e,q,t),d(e,ge,t),d(e,O,t),l(O,H);for(let u=0;u<y.length;u+=1)y[u]&&y[u].m(H,null);l(O,Le),l(O,L);for(let u=0;u<k.length;u+=1)k[u]&&k[u].m(L,null);P=!0},p(e,[t]){var Qe,ze;(!P||t&1)&&i!==(i=e[0].name+"")&&ve(v,i),(!P||t&1)&&J!==(J=e[0].name+"")&&ve(te,J);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const record = await pb.collection('${(Qe=e[0])==null?void 0:Qe.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final record = await pb.collection('${(ze=e[0])==null?void 0:ze.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `),F.$set(u),(!P||t&1)&&z!==(z=e[0].name+"")&&ve(oe,z),e[1]?g||(g=Ze(),g.c(),g.m($,null)):g&&(g.d(1),g=null),t&20&&(G=K(e[4]),y=Ge(y,t,Ve,1,e,G,He,H,ot,Ye,null,Xe)),t&20&&(U=K(e[4]),it(),k=Ge(k,t,je,1,e,U,Ue,L,rt,et,null,We),dt())},i(e){if(!P){j(F.$$.fragment,e),j(E.$$.fragment,e),j(A.$$.fragment,e);for(let t=0;t<U.length;t+=1)j(k[t]);P=!0}},o(e){V(F.$$.fragment,e),V(E.$$.fragment,e),V(A.$$.fragment,e);for(let t=0;t<k.length;t+=1)V(k[t]);P=!1},d(e){e&&(r(s),r(c),r(f),r(le),r(se),r(B),r(ne),r($),r(re),r(I),r(de),r(S),r(ce),r(x),r(pe),r(R),r(ke),r(q),r(ge),r(O)),W(F,e),g&&g.d(),W(E),W(A);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function bt(o,s,n){let i,v,{collection:p}=s,c=200,f=[];const w=C=>n(2,c=C.code);return o.$$set=C=>{"collection"in C&&n(0,p=C.collection)},o.$$.update=()=>{o.$$.dirty&1&&n(1,i=(p==null?void 0:p.viewRule)===null),o.$$.dirty&3&&p!=null&&p.id&&(f.push({code:200,body:JSON.stringify(Ke.dummyCollectionRecord(p),null,2)}),i&&f.push({code:403,body:`
                    {
                      "status": 403,
                      "message": "Only superusers can access this action.",
                      "data": {}
                    }
                `}),f.push({code:404,body:`
                {
                  "status": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},n(3,v=Ke.getApiExampleUrl(ct.baseURL)),[p,i,c,v,f,w]}class ht extends lt{constructor(s){super(),st(this,s,bt,ft,nt,{collection:0})}}export{ht as default};
