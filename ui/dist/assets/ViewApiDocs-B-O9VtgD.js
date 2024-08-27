import{S as lt,i as nt,s as st,N as tt,O as K,e as o,v as _,b as m,c as W,f as b,g as r,h as l,m as X,w as ve,P as Je,Q as ot,k as at,R as it,n as rt,t as Q,a as U,o as d,d as Y,C as Ke,A as dt,q as Z,r as ct}from"./index-D0DO79Dq.js";import{S as pt}from"./SdkTabs-DC6EUYpr.js";import{F as ut}from"./FieldsQueryParam-BwleQAus.js";function We(a,n,s){const i=a.slice();return i[6]=n[s],i}function Xe(a,n,s){const i=a.slice();return i[6]=n[s],i}function Ye(a){let n;return{c(){n=o("p"),n.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",b(n,"class","txt-hint txt-sm txt-right")},m(s,i){r(s,n,i)},d(s){s&&d(n)}}}function Ze(a,n){let s,i,v;function p(){return n[5](n[6])}return{key:a,first:null,c(){s=o("button"),s.textContent=`${n[6].code} `,b(s,"class","tab-item"),Z(s,"active",n[2]===n[6].code),this.first=s},m(c,f){r(c,s,f),i||(v=ct(s,"click",p),i=!0)},p(c,f){n=c,f&20&&Z(s,"active",n[2]===n[6].code)},d(c){c&&d(s),i=!1,v()}}}function et(a,n){let s,i,v,p;return i=new tt({props:{content:n[6].body}}),{key:a,first:null,c(){s=o("div"),W(i.$$.fragment),v=m(),b(s,"class","tab-item"),Z(s,"active",n[2]===n[6].code),this.first=s},m(c,f){r(c,s,f),X(i,s,null),l(s,v),p=!0},p(c,f){n=c,(!p||f&20)&&Z(s,"active",n[2]===n[6].code)},i(c){p||(Q(i.$$.fragment,c),p=!0)},o(c){U(i.$$.fragment,c),p=!1},d(c){c&&d(s),Y(i)}}}function ft(a){var je,Ve;let n,s,i=a[0].name+"",v,p,c,f,w,C,ee,j=a[0].name+"",te,$e,le,F,ne,S,se,$,V,ye,z,T,we,oe,G=a[0].name+"",ae,Ce,ie,Fe,re,B,de,q,ce,x,pe,R,ue,Re,I,O,fe,Oe,me,Pe,h,De,A,Te,Ae,Ee,be,Se,_e,Be,qe,xe,he,Ie,Me,E,ke,M,ge,P,H,y=[],He=new Map,Le,L,k=[],Ne=new Map,D;F=new pt({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        const record = await pb.collection('${(je=a[0])==null?void 0:je.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        final record = await pb.collection('${(Ve=a[0])==null?void 0:Ve.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let g=a[1]&&Ye();A=new tt({props:{content:"?expand=relField1,relField2.subRelField"}}),E=new ut({});let J=K(a[4]);const Qe=e=>e[6].code;for(let e=0;e<J.length;e+=1){let t=Xe(a,J,e),u=Qe(t);He.set(u,y[e]=Ze(u,t))}let N=K(a[4]);const Ue=e=>e[6].code;for(let e=0;e<N.length;e+=1){let t=We(a,N,e),u=Ue(t);Ne.set(u,k[e]=et(u,t))}return{c(){n=o("h3"),s=_("View ("),v=_(i),p=_(")"),c=m(),f=o("div"),w=o("p"),C=_("Fetch a single "),ee=o("strong"),te=_(j),$e=_(" record."),le=m(),W(F.$$.fragment),ne=m(),S=o("h6"),S.textContent="API details",se=m(),$=o("div"),V=o("strong"),V.textContent="GET",ye=m(),z=o("div"),T=o("p"),we=_("/api/collections/"),oe=o("strong"),ae=_(G),Ce=_("/records/"),ie=o("strong"),ie.textContent=":id",Fe=m(),g&&g.c(),re=m(),B=o("div"),B.textContent="Path Parameters",de=m(),q=o("table"),q.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to view.</td></tr></tbody>',ce=m(),x=o("div"),x.textContent="Query parameters",pe=m(),R=o("table"),ue=o("thead"),ue.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Re=m(),I=o("tbody"),O=o("tr"),fe=o("td"),fe.textContent="expand",Oe=m(),me=o("td"),me.innerHTML='<span class="label">String</span>',Pe=m(),h=o("td"),De=_(`Auto expand record relations. Ex.:
                `),W(A.$$.fragment),Te=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ae=o("br"),Ee=_(`
                The expanded relations will be appended to the record under the
                `),be=o("code"),be.textContent="expand",Se=_(" property (eg. "),_e=o("code"),_e.textContent='"expand": {"relField1": {...}, ...}',Be=_(`).
                `),qe=o("br"),xe=_(`
                Only the relations to which the request user has permissions to `),he=o("strong"),he.textContent="view",Ie=_(" will be expanded."),Me=m(),W(E.$$.fragment),ke=m(),M=o("div"),M.textContent="Responses",ge=m(),P=o("div"),H=o("div");for(let e=0;e<y.length;e+=1)y[e].c();Le=m(),L=o("div");for(let e=0;e<k.length;e+=1)k[e].c();b(n,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(S,"class","m-b-xs"),b(V,"class","label label-primary"),b(z,"class","content"),b($,"class","alert alert-info"),b(B,"class","section-title"),b(q,"class","table-compact table-border m-b-base"),b(x,"class","section-title"),b(R,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(H,"class","tabs-header compact combined left"),b(L,"class","tabs-content"),b(P,"class","tabs")},m(e,t){r(e,n,t),l(n,s),l(n,v),l(n,p),r(e,c,t),r(e,f,t),l(f,w),l(w,C),l(w,ee),l(ee,te),l(w,$e),r(e,le,t),X(F,e,t),r(e,ne,t),r(e,S,t),r(e,se,t),r(e,$,t),l($,V),l($,ye),l($,z),l(z,T),l(T,we),l(T,oe),l(oe,ae),l(T,Ce),l(T,ie),l($,Fe),g&&g.m($,null),r(e,re,t),r(e,B,t),r(e,de,t),r(e,q,t),r(e,ce,t),r(e,x,t),r(e,pe,t),r(e,R,t),l(R,ue),l(R,Re),l(R,I),l(I,O),l(O,fe),l(O,Oe),l(O,me),l(O,Pe),l(O,h),l(h,De),X(A,h,null),l(h,Te),l(h,Ae),l(h,Ee),l(h,be),l(h,Se),l(h,_e),l(h,Be),l(h,qe),l(h,xe),l(h,he),l(h,Ie),l(I,Me),X(E,I,null),r(e,ke,t),r(e,M,t),r(e,ge,t),r(e,P,t),l(P,H);for(let u=0;u<y.length;u+=1)y[u]&&y[u].m(H,null);l(P,Le),l(P,L);for(let u=0;u<k.length;u+=1)k[u]&&k[u].m(L,null);D=!0},p(e,[t]){var ze,Ge;(!D||t&1)&&i!==(i=e[0].name+"")&&ve(v,i),(!D||t&1)&&j!==(j=e[0].name+"")&&ve(te,j);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const record = await pb.collection('${(ze=e[0])==null?void 0:ze.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final record = await pb.collection('${(Ge=e[0])==null?void 0:Ge.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `),F.$set(u),(!D||t&1)&&G!==(G=e[0].name+"")&&ve(ae,G),e[1]?g||(g=Ye(),g.c(),g.m($,null)):g&&(g.d(1),g=null),t&20&&(J=K(e[4]),y=Je(y,t,Qe,1,e,J,He,H,ot,Ze,null,Xe)),t&20&&(N=K(e[4]),at(),k=Je(k,t,Ue,1,e,N,Ne,L,it,et,null,We),rt())},i(e){if(!D){Q(F.$$.fragment,e),Q(A.$$.fragment,e),Q(E.$$.fragment,e);for(let t=0;t<N.length;t+=1)Q(k[t]);D=!0}},o(e){U(F.$$.fragment,e),U(A.$$.fragment,e),U(E.$$.fragment,e);for(let t=0;t<k.length;t+=1)U(k[t]);D=!1},d(e){e&&(d(n),d(c),d(f),d(le),d(ne),d(S),d(se),d($),d(re),d(B),d(de),d(q),d(ce),d(x),d(pe),d(R),d(ke),d(M),d(ge),d(P)),Y(F,e),g&&g.d(),Y(A),Y(E);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function mt(a,n,s){let i,v,{collection:p}=n,c=200,f=[];const w=C=>s(2,c=C.code);return a.$$set=C=>{"collection"in C&&s(0,p=C.collection)},a.$$.update=()=>{a.$$.dirty&1&&s(1,i=(p==null?void 0:p.viewRule)===null),a.$$.dirty&3&&p!=null&&p.id&&(f.push({code:200,body:JSON.stringify(Ke.dummyCollectionRecord(p),null,2)}),i&&f.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}),f.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},s(3,v=Ke.getApiExampleUrl(dt.baseUrl)),[p,i,c,v,f,w]}class kt extends lt{constructor(n){super(),nt(this,n,mt,ft,st,{collection:0})}}export{kt as default};
