import{S as Re,i as Ee,s as Pe,V as Te,X as U,h as c,z as y,j as k,c as De,k as m,n as p,o as i,m as Ce,H as ee,Y as he,Z as Be,E as Oe,_ as Ie,G as qe,t as te,a as le,v as f,d as we,J as Ae,p as Me,l as z,q as He,W as Le}from"./index-CUlXN_4f.js";function ke(a,l,s){const n=a.slice();return n[6]=l[s],n}function ge(a,l,s){const n=a.slice();return n[6]=l[s],n}function ve(a){let l;return{c(){l=c("p"),l.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",m(l,"class","txt-hint txt-sm txt-right")},m(s,n){p(s,l,n)},d(s){s&&f(l)}}}function ye(a,l){let s,n,h;function r(){return l[5](l[6])}return{key:a,first:null,c(){s=c("button"),s.textContent=`${l[6].code} `,m(s,"class","tab-item"),z(s,"active",l[2]===l[6].code),this.first=s},m(o,d){p(o,s,d),n||(h=He(s,"click",r),n=!0)},p(o,d){l=o,d&20&&z(s,"active",l[2]===l[6].code)},d(o){o&&f(s),n=!1,h()}}}function $e(a,l){let s,n,h,r;return n=new Le({props:{content:l[6].body}}),{key:a,first:null,c(){s=c("div"),De(n.$$.fragment),h=k(),m(s,"class","tab-item"),z(s,"active",l[2]===l[6].code),this.first=s},m(o,d){p(o,s,d),Ce(n,s,null),i(s,h),r=!0},p(o,d){l=o,(!r||d&20)&&z(s,"active",l[2]===l[6].code)},i(o){r||(te(n.$$.fragment,o),r=!0)},o(o){le(n.$$.fragment,o),r=!1},d(o){o&&f(s),we(n)}}}function Se(a){var fe,me;let l,s,n=a[0].name+"",h,r,o,d,$,D,F,M=a[0].name+"",G,se,J,C,K,P,N,g,H,ae,L,E,ne,V,S=a[0].name+"",W,oe,X,ie,Y,T,Z,B,Q,O,x,w,I,v=[],ce=new Map,re,q,b=[],de=new Map,R;C=new Te({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(fe=a[0])==null?void 0:fe.name}').delete('RECORD_ID');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(me=a[0])==null?void 0:me.name}').delete('RECORD_ID');
    `}});let _=a[1]&&ve(),j=U(a[4]);const ue=e=>e[6].code;for(let e=0;e<j.length;e+=1){let t=ge(a,j,e),u=ue(t);ce.set(u,v[e]=ye(u,t))}let A=U(a[4]);const pe=e=>e[6].code;for(let e=0;e<A.length;e+=1){let t=ke(a,A,e),u=pe(t);de.set(u,b[e]=$e(u,t))}return{c(){l=c("h3"),s=y("Delete ("),h=y(n),r=y(")"),o=k(),d=c("div"),$=c("p"),D=y("Delete a single "),F=c("strong"),G=y(M),se=y(" record."),J=k(),De(C.$$.fragment),K=k(),P=c("h6"),P.textContent="API details",N=k(),g=c("div"),H=c("strong"),H.textContent="DELETE",ae=k(),L=c("div"),E=c("p"),ne=y("/api/collections/"),V=c("strong"),W=y(S),oe=y("/records/"),X=c("strong"),X.textContent=":id",ie=k(),_&&_.c(),Y=k(),T=c("div"),T.textContent="Path parameters",Z=k(),B=c("table"),B.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to delete.</td></tr></tbody>',Q=k(),O=c("div"),O.textContent="Responses",x=k(),w=c("div"),I=c("div");for(let e=0;e<v.length;e+=1)v[e].c();re=k(),q=c("div");for(let e=0;e<b.length;e+=1)b[e].c();m(l,"class","m-b-sm"),m(d,"class","content txt-lg m-b-sm"),m(P,"class","m-b-xs"),m(H,"class","label label-primary"),m(L,"class","content"),m(g,"class","alert alert-danger"),m(T,"class","section-title"),m(B,"class","table-compact table-border m-b-base"),m(O,"class","section-title"),m(I,"class","tabs-header compact combined left"),m(q,"class","tabs-content"),m(w,"class","tabs")},m(e,t){p(e,l,t),i(l,s),i(l,h),i(l,r),p(e,o,t),p(e,d,t),i(d,$),i($,D),i($,F),i(F,G),i($,se),p(e,J,t),Ce(C,e,t),p(e,K,t),p(e,P,t),p(e,N,t),p(e,g,t),i(g,H),i(g,ae),i(g,L),i(L,E),i(E,ne),i(E,V),i(V,W),i(E,oe),i(E,X),i(g,ie),_&&_.m(g,null),p(e,Y,t),p(e,T,t),p(e,Z,t),p(e,B,t),p(e,Q,t),p(e,O,t),p(e,x,t),p(e,w,t),i(w,I);for(let u=0;u<v.length;u+=1)v[u]&&v[u].m(I,null);i(w,re),i(w,q);for(let u=0;u<b.length;u+=1)b[u]&&b[u].m(q,null);R=!0},p(e,[t]){var be,_e;(!R||t&1)&&n!==(n=e[0].name+"")&&ee(h,n),(!R||t&1)&&M!==(M=e[0].name+"")&&ee(G,M);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').delete('RECORD_ID');
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').delete('RECORD_ID');
    `),C.$set(u),(!R||t&1)&&S!==(S=e[0].name+"")&&ee(W,S),e[1]?_||(_=ve(),_.c(),_.m(g,null)):_&&(_.d(1),_=null),t&20&&(j=U(e[4]),v=he(v,t,ue,1,e,j,ce,I,Be,ye,null,ge)),t&20&&(A=U(e[4]),Oe(),b=he(b,t,pe,1,e,A,de,q,Ie,$e,null,ke),qe())},i(e){if(!R){te(C.$$.fragment,e);for(let t=0;t<A.length;t+=1)te(b[t]);R=!0}},o(e){le(C.$$.fragment,e);for(let t=0;t<b.length;t+=1)le(b[t]);R=!1},d(e){e&&(f(l),f(o),f(d),f(J),f(K),f(P),f(N),f(g),f(Y),f(T),f(Z),f(B),f(Q),f(O),f(x),f(w)),we(C,e),_&&_.d();for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<b.length;t+=1)b[t].d()}}}function je(a,l,s){let n,h,{collection:r}=l,o=204,d=[];const $=D=>s(2,o=D.code);return a.$$set=D=>{"collection"in D&&s(0,r=D.collection)},a.$$.update=()=>{a.$$.dirty&1&&s(1,n=(r==null?void 0:r.deleteRule)===null),a.$$.dirty&3&&r!=null&&r.id&&(d.push({code:204,body:`
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
            `}))},s(3,h=Ae.getApiExampleUrl(Me.baseURL)),[r,n,o,h,d,$]}class ze extends Re{constructor(l){super(),Ee(this,l,je,Se,Pe,{collection:0})}}export{ze as default};