import{S as lt,i as nt,s as st,M as tt,N as K,e as a,w as _,b as m,c as W,f as b,g as r,h as l,m as X,x as ve,P as Je,Q as at,k as ot,R as it,n as rt,t as N,a as Q,o as d,d as Y,U as dt,C as Ke,p as ct,r as Z,u as pt}from"./index-a084d9d7.js";import{S as ut}from"./SdkTabs-ba0ec979.js";import{F as ft}from"./FieldsQueryParam-71e01e64.js";function We(o,n,s){const i=o.slice();return i[6]=n[s],i}function Xe(o,n,s){const i=o.slice();return i[6]=n[s],i}function Ye(o){let n;return{c(){n=a("p"),n.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",b(n,"class","txt-hint txt-sm txt-right")},m(s,i){r(s,n,i)},d(s){s&&d(n)}}}function Ze(o,n){let s,i,v;function p(){return n[5](n[6])}return{key:o,first:null,c(){s=a("button"),s.textContent=`${n[6].code} `,b(s,"class","tab-item"),Z(s,"active",n[2]===n[6].code),this.first=s},m(c,f){r(c,s,f),i||(v=pt(s,"click",p),i=!0)},p(c,f){n=c,f&20&&Z(s,"active",n[2]===n[6].code)},d(c){c&&d(s),i=!1,v()}}}function et(o,n){let s,i,v,p;return i=new tt({props:{content:n[6].body}}),{key:o,first:null,c(){s=a("div"),W(i.$$.fragment),v=m(),b(s,"class","tab-item"),Z(s,"active",n[2]===n[6].code),this.first=s},m(c,f){r(c,s,f),X(i,s,null),l(s,v),p=!0},p(c,f){n=c,(!p||f&20)&&Z(s,"active",n[2]===n[6].code)},i(c){p||(N(i.$$.fragment,c),p=!0)},o(c){Q(i.$$.fragment,c),p=!1},d(c){c&&d(s),Y(i)}}}function mt(o){var je,Ve;let n,s,i=o[0].name+"",v,p,c,f,y,C,ee,j=o[0].name+"",te,$e,le,F,ne,x,se,$,V,we,z,T,ye,ae,G=o[0].name+"",oe,Ce,ie,Fe,re,B,de,M,ce,A,pe,R,ue,Re,I,O,fe,Oe,me,Pe,h,De,E,Te,Ee,Se,be,xe,_e,Be,Me,Ae,he,Ie,qe,S,ke,q,ge,P,H,w=[],He=new Map,Le,L,k=[],Ue=new Map,D;F=new ut({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        const record = await pb.collection('${(je=o[0])==null?void 0:je.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        final record = await pb.collection('${(Ve=o[0])==null?void 0:Ve.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let g=o[1]&&Ye();E=new tt({props:{content:"?expand=relField1,relField2.subRelField"}}),S=new ft({});let J=K(o[4]);const Ne=e=>e[6].code;for(let e=0;e<J.length;e+=1){let t=Xe(o,J,e),u=Ne(t);He.set(u,w[e]=Ze(u,t))}let U=K(o[4]);const Qe=e=>e[6].code;for(let e=0;e<U.length;e+=1){let t=We(o,U,e),u=Qe(t);Ue.set(u,k[e]=et(u,t))}return{c(){n=a("h3"),s=_("View ("),v=_(i),p=_(")"),c=m(),f=a("div"),y=a("p"),C=_("Fetch a single "),ee=a("strong"),te=_(j),$e=_(" record."),le=m(),W(F.$$.fragment),ne=m(),x=a("h6"),x.textContent="API details",se=m(),$=a("div"),V=a("strong"),V.textContent="GET",we=m(),z=a("div"),T=a("p"),ye=_("/api/collections/"),ae=a("strong"),oe=_(G),Ce=_("/records/"),ie=a("strong"),ie.textContent=":id",Fe=m(),g&&g.c(),re=m(),B=a("div"),B.textContent="Path Parameters",de=m(),M=a("table"),M.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to view.</td></tr></tbody>',ce=m(),A=a("div"),A.textContent="Query parameters",pe=m(),R=a("table"),ue=a("thead"),ue.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Re=m(),I=a("tbody"),O=a("tr"),fe=a("td"),fe.textContent="expand",Oe=m(),me=a("td"),me.innerHTML='<span class="label">String</span>',Pe=m(),h=a("td"),De=_(`Auto expand record relations. Ex.:
                `),W(E.$$.fragment),Te=_(`
                Supports up to 6-levels depth nested relations expansion. `),Ee=a("br"),Se=_(`
                The expanded relations will be appended to the record under the
                `),be=a("code"),be.textContent="expand",xe=_(" property (eg. "),_e=a("code"),_e.textContent='"expand": {"relField1": {...}, ...}',Be=_(`).
                `),Me=a("br"),Ae=_(`
                Only the relations to which the request user has permissions to `),he=a("strong"),he.textContent="view",Ie=_(" will be expanded."),qe=m(),W(S.$$.fragment),ke=m(),q=a("div"),q.textContent="Responses",ge=m(),P=a("div"),H=a("div");for(let e=0;e<w.length;e+=1)w[e].c();Le=m(),L=a("div");for(let e=0;e<k.length;e+=1)k[e].c();b(n,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(x,"class","m-b-xs"),b(V,"class","label label-primary"),b(z,"class","content"),b($,"class","alert alert-info"),b(B,"class","section-title"),b(M,"class","table-compact table-border m-b-base"),b(A,"class","section-title"),b(R,"class","table-compact table-border m-b-base"),b(q,"class","section-title"),b(H,"class","tabs-header compact left"),b(L,"class","tabs-content"),b(P,"class","tabs")},m(e,t){r(e,n,t),l(n,s),l(n,v),l(n,p),r(e,c,t),r(e,f,t),l(f,y),l(y,C),l(y,ee),l(ee,te),l(y,$e),r(e,le,t),X(F,e,t),r(e,ne,t),r(e,x,t),r(e,se,t),r(e,$,t),l($,V),l($,we),l($,z),l(z,T),l(T,ye),l(T,ae),l(ae,oe),l(T,Ce),l(T,ie),l($,Fe),g&&g.m($,null),r(e,re,t),r(e,B,t),r(e,de,t),r(e,M,t),r(e,ce,t),r(e,A,t),r(e,pe,t),r(e,R,t),l(R,ue),l(R,Re),l(R,I),l(I,O),l(O,fe),l(O,Oe),l(O,me),l(O,Pe),l(O,h),l(h,De),X(E,h,null),l(h,Te),l(h,Ee),l(h,Se),l(h,be),l(h,xe),l(h,_e),l(h,Be),l(h,Me),l(h,Ae),l(h,he),l(h,Ie),l(I,qe),X(S,I,null),r(e,ke,t),r(e,q,t),r(e,ge,t),r(e,P,t),l(P,H);for(let u=0;u<w.length;u+=1)w[u]&&w[u].m(H,null);l(P,Le),l(P,L);for(let u=0;u<k.length;u+=1)k[u]&&k[u].m(L,null);D=!0},p(e,[t]){var ze,Ge;(!D||t&1)&&i!==(i=e[0].name+"")&&ve(v,i),(!D||t&1)&&j!==(j=e[0].name+"")&&ve(te,j);const u={};t&9&&(u.js=`
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
    `),F.$set(u),(!D||t&1)&&G!==(G=e[0].name+"")&&ve(oe,G),e[1]?g||(g=Ye(),g.c(),g.m($,null)):g&&(g.d(1),g=null),t&20&&(J=K(e[4]),w=Je(w,t,Ne,1,e,J,He,H,at,Ze,null,Xe)),t&20&&(U=K(e[4]),ot(),k=Je(k,t,Qe,1,e,U,Ue,L,it,et,null,We),rt())},i(e){if(!D){N(F.$$.fragment,e),N(E.$$.fragment,e),N(S.$$.fragment,e);for(let t=0;t<U.length;t+=1)N(k[t]);D=!0}},o(e){Q(F.$$.fragment,e),Q(E.$$.fragment,e),Q(S.$$.fragment,e);for(let t=0;t<k.length;t+=1)Q(k[t]);D=!1},d(e){e&&(d(n),d(c),d(f),d(le),d(ne),d(x),d(se),d($),d(re),d(B),d(de),d(M),d(ce),d(A),d(pe),d(R),d(ke),d(q),d(ge),d(P)),Y(F,e),g&&g.d(),Y(E),Y(S);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function bt(o,n,s){let i,v,{collection:p=new dt}=n,c=200,f=[];const y=C=>s(2,c=C.code);return o.$$set=C=>{"collection"in C&&s(0,p=C.collection)},o.$$.update=()=>{o.$$.dirty&1&&s(1,i=(p==null?void 0:p.viewRule)===null),o.$$.dirty&3&&p!=null&&p.id&&(f.push({code:200,body:JSON.stringify(Ke.dummyCollectionRecord(p),null,2)}),i&&f.push({code:403,body:`
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
            `}))},s(3,v=Ke.getApiExampleUrl(ct.baseUrl)),[p,i,c,v,f,y]}class gt extends lt{constructor(n){super(),nt(this,n,bt,mt,st,{collection:0})}}export{gt as default};
