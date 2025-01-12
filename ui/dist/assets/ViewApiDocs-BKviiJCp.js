import{S as lt,i as st,s as nt,V as ot,W as tt,X as K,h as o,z as _,j as b,c as W,k as m,n as r,o as l,m as X,H as ve,Y as Qe,Z as at,E as it,_ as rt,G as dt,t as U,a as V,v as d,d as Y,J as Ke,p as ct,l as Z,q as pt}from"./index-SKn09NMF.js";import{F as ut}from"./FieldsQueryParam-NXRpAlxi.js";function We(a,s,n){const i=a.slice();return i[6]=s[n],i}function Xe(a,s,n){const i=a.slice();return i[6]=s[n],i}function Ye(a){let s;return{c(){s=o("p"),s.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",m(s,"class","txt-hint txt-sm txt-right")},m(n,i){r(n,s,i)},d(n){n&&d(s)}}}function Ze(a,s){let n,i,v;function p(){return s[5](s[6])}return{key:a,first:null,c(){n=o("button"),n.textContent=`${s[6].code} `,m(n,"class","tab-item"),Z(n,"active",s[2]===s[6].code),this.first=n},m(c,f){r(c,n,f),i||(v=pt(n,"click",p),i=!0)},p(c,f){s=c,f&20&&Z(n,"active",s[2]===s[6].code)},d(c){c&&d(n),i=!1,v()}}}function et(a,s){let n,i,v,p;return i=new tt({props:{content:s[6].body}}),{key:a,first:null,c(){n=o("div"),W(i.$$.fragment),v=b(),m(n,"class","tab-item"),Z(n,"active",s[2]===s[6].code),this.first=n},m(c,f){r(c,n,f),X(i,n,null),l(n,v),p=!0},p(c,f){s=c,(!p||f&20)&&Z(n,"active",s[2]===s[6].code)},i(c){p||(U(i.$$.fragment,c),p=!0)},o(c){V(i.$$.fragment,c),p=!1},d(c){c&&d(n),Y(i)}}}function ft(a){var ze,Ge;let s,n,i=a[0].name+"",v,p,c,f,w,C,ee,z=a[0].name+"",te,$e,le,F,se,S,ne,$,G,ye,J,T,we,oe,N=a[0].name+"",ae,Ce,ie,Fe,re,q,de,x,ce,A,pe,R,ue,Re,H,O,fe,Oe,be,De,h,Pe,E,Te,Ee,Be,me,Se,_e,qe,xe,Ae,he,He,Ie,B,ke,I,ge,D,M,y=[],Me=new Map,Le,L,k=[],je=new Map,P;F=new ot({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        const record = await pb.collection('${(ze=a[0])==null?void 0:ze.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        final record = await pb.collection('${(Ge=a[0])==null?void 0:Ge.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let g=a[1]&&Ye();E=new tt({props:{content:"?expand=relField1,relField2.subRelField"}}),B=new ut({});let Q=K(a[4]);const Ue=e=>e[6].code;for(let e=0;e<Q.length;e+=1){let t=Xe(a,Q,e),u=Ue(t);Me.set(u,y[e]=Ze(u,t))}let j=K(a[4]);const Ve=e=>e[6].code;for(let e=0;e<j.length;e+=1){let t=We(a,j,e),u=Ve(t);je.set(u,k[e]=et(u,t))}return{c(){s=o("h3"),n=_("View ("),v=_(i),p=_(")"),c=b(),f=o("div"),w=o("p"),C=_("Fetch a single "),ee=o("strong"),te=_(z),$e=_(" record."),le=b(),W(F.$$.fragment),se=b(),S=o("h6"),S.textContent="API details",ne=b(),$=o("div"),G=o("strong"),G.textContent="GET",ye=b(),J=o("div"),T=o("p"),we=_("/api/collections/"),oe=o("strong"),ae=_(N),Ce=_("/records/"),ie=o("strong"),ie.textContent=":id",Fe=b(),g&&g.c(),re=b(),q=o("div"),q.textContent="Path Parameters",de=b(),x=o("table"),x.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to view.</td></tr></tbody>',ce=b(),A=o("div"),A.textContent="Query parameters",pe=b(),R=o("table"),ue=o("thead"),ue.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Re=b(),H=o("tbody"),O=o("tr"),fe=o("td"),fe.textContent="expand",Oe=b(),be=o("td"),be.innerHTML='<span class="label">String</span>',De=b(),h=o("td"),Pe=_(`Auto expand record relations. Ex.:
                `),W(E.$$.fragment),Te=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ee=o("br"),Be=_(`
                The expanded relations will be appended to the record under the
                `),me=o("code"),me.textContent="expand",Se=_(" property (eg. "),_e=o("code"),_e.textContent='"expand": {"relField1": {...}, ...}',qe=_(`).
                `),xe=o("br"),Ae=_(`
                Only the relations to which the request user has permissions to `),he=o("strong"),he.textContent="view",He=_(" will be expanded."),Ie=b(),W(B.$$.fragment),ke=b(),I=o("div"),I.textContent="Responses",ge=b(),D=o("div"),M=o("div");for(let e=0;e<y.length;e+=1)y[e].c();Le=b(),L=o("div");for(let e=0;e<k.length;e+=1)k[e].c();m(s,"class","m-b-sm"),m(f,"class","content txt-lg m-b-sm"),m(S,"class","m-b-xs"),m(G,"class","label label-primary"),m(J,"class","content"),m($,"class","alert alert-info"),m(q,"class","section-title"),m(x,"class","table-compact table-border m-b-base"),m(A,"class","section-title"),m(R,"class","table-compact table-border m-b-base"),m(I,"class","section-title"),m(M,"class","tabs-header compact combined left"),m(L,"class","tabs-content"),m(D,"class","tabs")},m(e,t){r(e,s,t),l(s,n),l(s,v),l(s,p),r(e,c,t),r(e,f,t),l(f,w),l(w,C),l(w,ee),l(ee,te),l(w,$e),r(e,le,t),X(F,e,t),r(e,se,t),r(e,S,t),r(e,ne,t),r(e,$,t),l($,G),l($,ye),l($,J),l(J,T),l(T,we),l(T,oe),l(oe,ae),l(T,Ce),l(T,ie),l($,Fe),g&&g.m($,null),r(e,re,t),r(e,q,t),r(e,de,t),r(e,x,t),r(e,ce,t),r(e,A,t),r(e,pe,t),r(e,R,t),l(R,ue),l(R,Re),l(R,H),l(H,O),l(O,fe),l(O,Oe),l(O,be),l(O,De),l(O,h),l(h,Pe),X(E,h,null),l(h,Te),l(h,Ee),l(h,Be),l(h,me),l(h,Se),l(h,_e),l(h,qe),l(h,xe),l(h,Ae),l(h,he),l(h,He),l(H,Ie),X(B,H,null),r(e,ke,t),r(e,I,t),r(e,ge,t),r(e,D,t),l(D,M);for(let u=0;u<y.length;u+=1)y[u]&&y[u].m(M,null);l(D,Le),l(D,L);for(let u=0;u<k.length;u+=1)k[u]&&k[u].m(L,null);P=!0},p(e,[t]){var Je,Ne;(!P||t&1)&&i!==(i=e[0].name+"")&&ve(v,i),(!P||t&1)&&z!==(z=e[0].name+"")&&ve(te,z);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const record = await pb.collection('${(Je=e[0])==null?void 0:Je.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final record = await pb.collection('${(Ne=e[0])==null?void 0:Ne.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `),F.$set(u),(!P||t&1)&&N!==(N=e[0].name+"")&&ve(ae,N),e[1]?g||(g=Ye(),g.c(),g.m($,null)):g&&(g.d(1),g=null),t&20&&(Q=K(e[4]),y=Qe(y,t,Ue,1,e,Q,Me,M,at,Ze,null,Xe)),t&20&&(j=K(e[4]),it(),k=Qe(k,t,Ve,1,e,j,je,L,rt,et,null,We),dt())},i(e){if(!P){U(F.$$.fragment,e),U(E.$$.fragment,e),U(B.$$.fragment,e);for(let t=0;t<j.length;t+=1)U(k[t]);P=!0}},o(e){V(F.$$.fragment,e),V(E.$$.fragment,e),V(B.$$.fragment,e);for(let t=0;t<k.length;t+=1)V(k[t]);P=!1},d(e){e&&(d(s),d(c),d(f),d(le),d(se),d(S),d(ne),d($),d(re),d(q),d(de),d(x),d(ce),d(A),d(pe),d(R),d(ke),d(I),d(ge),d(D)),Y(F,e),g&&g.d(),Y(E),Y(B);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function bt(a,s,n){let i,v,{collection:p}=s,c=200,f=[];const w=C=>n(2,c=C.code);return a.$$set=C=>{"collection"in C&&n(0,p=C.collection)},a.$$.update=()=>{a.$$.dirty&1&&n(1,i=(p==null?void 0:p.viewRule)===null),a.$$.dirty&3&&p!=null&&p.id&&(f.push({code:200,body:JSON.stringify(Ke.dummyCollectionRecord(p),null,2)}),i&&f.push({code:403,body:`
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
