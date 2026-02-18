import{S as Re,i as Ee,s as Pe,V as Te,X as j,h as p,d as De,t as te,a as le,I as ee,Z as he,_ as Be,C as Ie,$ as Oe,D as Ae,l as f,n as i,m as Ce,u as c,A as $,v as k,c as we,w as m,J as Me,p as qe,k as z,o as Le,W as Se}from"./index-TiFsHbkW.js";function ke(a,l,s){const n=a.slice();return n[6]=l[s],n}function ge(a,l,s){const n=a.slice();return n[6]=l[s],n}function ve(a){let l;return{c(){l=c("p"),l.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",m(l,"class","txt-hint txt-sm txt-right")},m(s,n){f(s,l,n)},d(s){s&&p(l)}}}function $e(a,l){let s,n,h;function r(){return l[5](l[6])}return{key:a,first:null,c(){s=c("button"),s.textContent=`${l[6].code} `,m(s,"class","tab-item"),z(s,"active",l[2]===l[6].code),this.first=s},m(o,d){f(o,s,d),n||(h=Le(s,"click",r),n=!0)},p(o,d){l=o,d&20&&z(s,"active",l[2]===l[6].code)},d(o){o&&p(s),n=!1,h()}}}function ye(a,l){let s,n,h,r;return n=new Se({props:{content:l[6].body}}),{key:a,first:null,c(){s=c("div"),we(n.$$.fragment),h=k(),m(s,"class","tab-item"),z(s,"active",l[2]===l[6].code),this.first=s},m(o,d){f(o,s,d),Ce(n,s,null),i(s,h),r=!0},p(o,d){l=o,(!r||d&20)&&z(s,"active",l[2]===l[6].code)},i(o){r||(le(n.$$.fragment,o),r=!0)},o(o){te(n.$$.fragment,o),r=!1},d(o){o&&p(s),De(n)}}}function He(a){var fe,me;let l,s,n=a[0].name+"",h,r,o,d,y,D,F,q=a[0].name+"",J,se,K,C,N,P,V,g,L,ae,S,E,ne,W,H=a[0].name+"",X,oe,Z,ie,G,T,Q,B,Y,I,x,w,O,v=[],ce=new Map,re,A,b=[],de=new Map,R;C=new Te({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(fe=a[0])==null?void 0:fe.name}').delete('RECORD_ID');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(me=a[0])==null?void 0:me.name}').delete('RECORD_ID');
    `}});let _=a[1]&&ve(),U=j(a[4]);const ue=e=>e[6].code;for(let e=0;e<U.length;e+=1){let t=ge(a,U,e),u=ue(t);ce.set(u,v[e]=$e(u,t))}let M=j(a[4]);const pe=e=>e[6].code;for(let e=0;e<M.length;e+=1){let t=ke(a,M,e),u=pe(t);de.set(u,b[e]=ye(u,t))}return{c(){l=c("h3"),s=$("Delete ("),h=$(n),r=$(")"),o=k(),d=c("div"),y=c("p"),D=$("Delete a single "),F=c("strong"),J=$(q),se=$(" record."),K=k(),we(C.$$.fragment),N=k(),P=c("h6"),P.textContent="API details",V=k(),g=c("div"),L=c("strong"),L.textContent="DELETE",ae=k(),S=c("div"),E=c("p"),ne=$("/api/collections/"),W=c("strong"),X=$(H),oe=$("/records/"),Z=c("strong"),Z.textContent=":id",ie=k(),_&&_.c(),G=k(),T=c("div"),T.textContent="Path parameters",Q=k(),B=c("table"),B.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to delete.</td></tr></tbody>',Y=k(),I=c("div"),I.textContent="Responses",x=k(),w=c("div"),O=c("div");for(let e=0;e<v.length;e+=1)v[e].c();re=k(),A=c("div");for(let e=0;e<b.length;e+=1)b[e].c();m(l,"class","m-b-sm"),m(d,"class","content txt-lg m-b-sm"),m(P,"class","m-b-xs"),m(L,"class","label label-primary"),m(S,"class","content"),m(g,"class","alert alert-danger"),m(T,"class","section-title"),m(B,"class","table-compact table-border m-b-base"),m(I,"class","section-title"),m(O,"class","tabs-header compact combined left"),m(A,"class","tabs-content"),m(w,"class","tabs")},m(e,t){f(e,l,t),i(l,s),i(l,h),i(l,r),f(e,o,t),f(e,d,t),i(d,y),i(y,D),i(y,F),i(F,J),i(y,se),f(e,K,t),Ce(C,e,t),f(e,N,t),f(e,P,t),f(e,V,t),f(e,g,t),i(g,L),i(g,ae),i(g,S),i(S,E),i(E,ne),i(E,W),i(W,X),i(E,oe),i(E,Z),i(g,ie),_&&_.m(g,null),f(e,G,t),f(e,T,t),f(e,Q,t),f(e,B,t),f(e,Y,t),f(e,I,t),f(e,x,t),f(e,w,t),i(w,O);for(let u=0;u<v.length;u+=1)v[u]&&v[u].m(O,null);i(w,re),i(w,A);for(let u=0;u<b.length;u+=1)b[u]&&b[u].m(A,null);R=!0},p(e,[t]){var be,_e;(!R||t&1)&&n!==(n=e[0].name+"")&&ee(h,n),(!R||t&1)&&q!==(q=e[0].name+"")&&ee(J,q);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').delete('RECORD_ID');
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').delete('RECORD_ID');
    `),C.$set(u),(!R||t&1)&&H!==(H=e[0].name+"")&&ee(X,H),e[1]?_||(_=ve(),_.c(),_.m(g,null)):_&&(_.d(1),_=null),t&20&&(U=j(e[4]),v=he(v,t,ue,1,e,U,ce,O,Be,$e,null,ge)),t&20&&(M=j(e[4]),Ie(),b=he(b,t,pe,1,e,M,de,A,Oe,ye,null,ke),Ae())},i(e){if(!R){le(C.$$.fragment,e);for(let t=0;t<M.length;t+=1)le(b[t]);R=!0}},o(e){te(C.$$.fragment,e);for(let t=0;t<b.length;t+=1)te(b[t]);R=!1},d(e){e&&(p(l),p(o),p(d),p(K),p(N),p(P),p(V),p(g),p(G),p(T),p(Q),p(B),p(Y),p(I),p(x),p(w)),De(C,e),_&&_.d();for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<b.length;t+=1)b[t].d()}}}function Ue(a,l,s){let n,h,{collection:r}=l,o=204,d=[];const y=D=>s(2,o=D.code);return a.$$set=D=>{"collection"in D&&s(0,r=D.collection)},a.$$.update=()=>{a.$$.dirty&1&&s(1,n=(r==null?void 0:r.deleteRule)===null),a.$$.dirty&3&&r!=null&&r.id&&(d.push({code:204,body:`
                null
            `}),d.push({code:400,body:`
                {
                  "status": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `}),n&&d.push({code:403,body:`
                    {
                      "status": 403,
                      "message": "Only superusers can access this action.",
                      "data": {}
                    }
                `}),d.push({code:404,body:`
                {
                  "status": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},s(3,h=Me.getApiExampleUrl(qe.baseURL)),[r,n,o,h,d,y]}class ze extends Re{constructor(l){super(),Ee(this,l,Ue,He,Pe,{collection:0})}}export{ze as default};
