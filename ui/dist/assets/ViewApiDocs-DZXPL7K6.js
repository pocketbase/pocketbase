import{S as lt,i as st,s as nt,U as ot,V as tt,W as K,f as o,y as _,h as b,c as W,j as m,l as r,n as l,m as X,G as ve,X as Je,Y as at,D as it,Z as rt,E as dt,t as j,a as V,u as d,d as Y,I as Ke,p as ct,k as Z,o as pt}from"./index-B224lkEB.js";import{F as ut}from"./FieldsQueryParam-CxM6Mw06.js";function We(a,s,n){const i=a.slice();return i[6]=s[n],i}function Xe(a,s,n){const i=a.slice();return i[6]=s[n],i}function Ye(a){let s;return{c(){s=o("p"),s.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",m(s,"class","txt-hint txt-sm txt-right")},m(n,i){r(n,s,i)},d(n){n&&d(s)}}}function Ze(a,s){let n,i,v;function p(){return s[5](s[6])}return{key:a,first:null,c(){n=o("button"),n.textContent=`${s[6].code} `,m(n,"class","tab-item"),Z(n,"active",s[2]===s[6].code),this.first=n},m(c,f){r(c,n,f),i||(v=pt(n,"click",p),i=!0)},p(c,f){s=c,f&20&&Z(n,"active",s[2]===s[6].code)},d(c){c&&d(n),i=!1,v()}}}function et(a,s){let n,i,v,p;return i=new tt({props:{content:s[6].body}}),{key:a,first:null,c(){n=o("div"),W(i.$$.fragment),v=b(),m(n,"class","tab-item"),Z(n,"active",s[2]===s[6].code),this.first=n},m(c,f){r(c,n,f),X(i,n,null),l(n,v),p=!0},p(c,f){s=c,(!p||f&20)&&Z(n,"active",s[2]===s[6].code)},i(c){p||(j(i.$$.fragment,c),p=!0)},o(c){V(i.$$.fragment,c),p=!1},d(c){c&&d(n),Y(i)}}}function ft(a){var Ge,Ne;let s,n,i=a[0].name+"",v,p,c,f,w,C,ee,G=a[0].name+"",te,$e,le,F,se,I,ne,$,N,ye,Q,T,we,oe,z=a[0].name+"",ae,Ce,ie,Fe,re,S,de,x,ce,A,pe,R,ue,Re,M,D,fe,De,be,Oe,h,Pe,E,Te,Ee,Be,me,Ie,_e,Se,xe,Ae,he,Me,qe,B,ke,q,ge,O,H,y=[],He=new Map,Le,L,k=[],Ue=new Map,P;F=new ot({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        const record = await pb.collection('${(Ge=a[0])==null?void 0:Ge.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        final record = await pb.collection('${(Ne=a[0])==null?void 0:Ne.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let g=a[1]&&Ye();E=new tt({props:{content:"?expand=relField1,relField2.subRelField"}}),B=new ut({});let J=K(a[4]);const je=e=>e[6].code;for(let e=0;e<J.length;e+=1){let t=Xe(a,J,e),u=je(t);He.set(u,y[e]=Ze(u,t))}let U=K(a[4]);const Ve=e=>e[6].code;for(let e=0;e<U.length;e+=1){let t=We(a,U,e),u=Ve(t);Ue.set(u,k[e]=et(u,t))}return{c(){s=o("h3"),n=_("View ("),v=_(i),p=_(")"),c=b(),f=o("div"),w=o("p"),C=_("Fetch a single "),ee=o("strong"),te=_(G),$e=_(" record."),le=b(),W(F.$$.fragment),se=b(),I=o("h6"),I.textContent="API details",ne=b(),$=o("div"),N=o("strong"),N.textContent="GET",ye=b(),Q=o("div"),T=o("p"),we=_("/api/collections/"),oe=o("strong"),ae=_(z),Ce=_("/records/"),ie=o("strong"),ie.textContent=":id",Fe=b(),g&&g.c(),re=b(),S=o("div"),S.textContent="Path Parameters",de=b(),x=o("table"),x.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to view.</td></tr></tbody>',ce=b(),A=o("div"),A.textContent="Query parameters",pe=b(),R=o("table"),ue=o("thead"),ue.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Re=b(),M=o("tbody"),D=o("tr"),fe=o("td"),fe.textContent="expand",De=b(),be=o("td"),be.innerHTML='<span class="label">String</span>',Oe=b(),h=o("td"),Pe=_(`Auto expand record relations. Ex.:
                `),W(E.$$.fragment),Te=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ee=o("br"),Be=_(`
                The expanded relations will be appended to the record under the
                `),me=o("code"),me.textContent="expand",Ie=_(" property (eg. "),_e=o("code"),_e.textContent='"expand": {"relField1": {...}, ...}',Se=_(`).
                `),xe=o("br"),Ae=_(`
                Only the relations to which the request user has permissions to `),he=o("strong"),he.textContent="view",Me=_(" will be expanded."),qe=b(),W(B.$$.fragment),ke=b(),q=o("div"),q.textContent="Responses",ge=b(),O=o("div"),H=o("div");for(let e=0;e<y.length;e+=1)y[e].c();Le=b(),L=o("div");for(let e=0;e<k.length;e+=1)k[e].c();m(s,"class","m-b-sm"),m(f,"class","content txt-lg m-b-sm"),m(I,"class","m-b-xs"),m(N,"class","label label-primary"),m(Q,"class","content"),m($,"class","alert alert-info"),m(S,"class","section-title"),m(x,"class","table-compact table-border m-b-base"),m(A,"class","section-title"),m(R,"class","table-compact table-border m-b-base"),m(q,"class","section-title"),m(H,"class","tabs-header compact combined left"),m(L,"class","tabs-content"),m(O,"class","tabs")},m(e,t){r(e,s,t),l(s,n),l(s,v),l(s,p),r(e,c,t),r(e,f,t),l(f,w),l(w,C),l(w,ee),l(ee,te),l(w,$e),r(e,le,t),X(F,e,t),r(e,se,t),r(e,I,t),r(e,ne,t),r(e,$,t),l($,N),l($,ye),l($,Q),l(Q,T),l(T,we),l(T,oe),l(oe,ae),l(T,Ce),l(T,ie),l($,Fe),g&&g.m($,null),r(e,re,t),r(e,S,t),r(e,de,t),r(e,x,t),r(e,ce,t),r(e,A,t),r(e,pe,t),r(e,R,t),l(R,ue),l(R,Re),l(R,M),l(M,D),l(D,fe),l(D,De),l(D,be),l(D,Oe),l(D,h),l(h,Pe),X(E,h,null),l(h,Te),l(h,Ee),l(h,Be),l(h,me),l(h,Ie),l(h,_e),l(h,Se),l(h,xe),l(h,Ae),l(h,he),l(h,Me),l(M,qe),X(B,M,null),r(e,ke,t),r(e,q,t),r(e,ge,t),r(e,O,t),l(O,H);for(let u=0;u<y.length;u+=1)y[u]&&y[u].m(H,null);l(O,Le),l(O,L);for(let u=0;u<k.length;u+=1)k[u]&&k[u].m(L,null);P=!0},p(e,[t]){var Qe,ze;(!P||t&1)&&i!==(i=e[0].name+"")&&ve(v,i),(!P||t&1)&&G!==(G=e[0].name+"")&&ve(te,G);const u={};t&9&&(u.js=`
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
    `),F.$set(u),(!P||t&1)&&z!==(z=e[0].name+"")&&ve(ae,z),e[1]?g||(g=Ye(),g.c(),g.m($,null)):g&&(g.d(1),g=null),t&20&&(J=K(e[4]),y=Je(y,t,je,1,e,J,He,H,at,Ze,null,Xe)),t&20&&(U=K(e[4]),it(),k=Je(k,t,Ve,1,e,U,Ue,L,rt,et,null,We),dt())},i(e){if(!P){j(F.$$.fragment,e),j(E.$$.fragment,e),j(B.$$.fragment,e);for(let t=0;t<U.length;t+=1)j(k[t]);P=!0}},o(e){V(F.$$.fragment,e),V(E.$$.fragment,e),V(B.$$.fragment,e);for(let t=0;t<k.length;t+=1)V(k[t]);P=!1},d(e){e&&(d(s),d(c),d(f),d(le),d(se),d(I),d(ne),d($),d(re),d(S),d(de),d(x),d(ce),d(A),d(pe),d(R),d(ke),d(q),d(ge),d(O)),Y(F,e),g&&g.d(),Y(E),Y(B);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function bt(a,s,n){let i,v,{collection:p}=s,c=200,f=[];const w=C=>n(2,c=C.code);return a.$$set=C=>{"collection"in C&&n(0,p=C.collection)},a.$$.update=()=>{a.$$.dirty&1&&n(1,i=(p==null?void 0:p.viewRule)===null),a.$$.dirty&3&&p!=null&&p.id&&(f.push({code:200,body:JSON.stringify(Ke.dummyCollectionRecord(p),null,2)}),i&&f.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only superusers can access this action.",
                      "data": {}
                    }
                `}),f.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},n(3,v=Ke.getApiExampleUrl(ct.baseURL)),[p,i,c,v,f,w]}class ht extends lt{constructor(s){super(),st(this,s,bt,ft,nt,{collection:0})}}export{ht as default};
