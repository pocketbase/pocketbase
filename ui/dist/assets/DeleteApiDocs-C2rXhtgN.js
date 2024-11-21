import{S as Re,i as Ee,s as Pe,U as Te,W as H,f as c,y,h as k,c as $e,j as m,l as p,n as i,m as Ce,G as ee,X as he,Y as Be,D as Ie,Z as Oe,E as Ae,t as te,a as le,u as f,d as we,I as Me,p as qe,k as z,o as Le,V as Se}from"./index-BigBIDNp.js";function ke(a,l,s){const o=a.slice();return o[6]=l[s],o}function ge(a,l,s){const o=a.slice();return o[6]=l[s],o}function ve(a){let l;return{c(){l=c("p"),l.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",m(l,"class","txt-hint txt-sm txt-right")},m(s,o){p(s,l,o)},d(s){s&&f(l)}}}function ye(a,l){let s,o,h;function r(){return l[5](l[6])}return{key:a,first:null,c(){s=c("button"),s.textContent=`${l[6].code} `,m(s,"class","tab-item"),z(s,"active",l[2]===l[6].code),this.first=s},m(n,d){p(n,s,d),o||(h=Le(s,"click",r),o=!0)},p(n,d){l=n,d&20&&z(s,"active",l[2]===l[6].code)},d(n){n&&f(s),o=!1,h()}}}function De(a,l){let s,o,h,r;return o=new Se({props:{content:l[6].body}}),{key:a,first:null,c(){s=c("div"),$e(o.$$.fragment),h=k(),m(s,"class","tab-item"),z(s,"active",l[2]===l[6].code),this.first=s},m(n,d){p(n,s,d),Ce(o,s,null),i(s,h),r=!0},p(n,d){l=n,(!r||d&20)&&z(s,"active",l[2]===l[6].code)},i(n){r||(te(o.$$.fragment,n),r=!0)},o(n){le(o.$$.fragment,n),r=!1},d(n){n&&f(s),we(o)}}}function Ue(a){var fe,me;let l,s,o=a[0].name+"",h,r,n,d,D,$,F,q=a[0].name+"",G,se,K,C,N,P,V,g,L,ae,S,E,oe,W,U=a[0].name+"",X,ne,Y,ie,Z,T,J,B,Q,I,x,w,O,v=[],ce=new Map,re,A,b=[],de=new Map,R;C=new Te({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(fe=a[0])==null?void 0:fe.name}').delete('RECORD_ID');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(me=a[0])==null?void 0:me.name}').delete('RECORD_ID');
    `}});let _=a[1]&&ve(),j=H(a[4]);const ue=e=>e[6].code;for(let e=0;e<j.length;e+=1){let t=ge(a,j,e),u=ue(t);ce.set(u,v[e]=ye(u,t))}let M=H(a[4]);const pe=e=>e[6].code;for(let e=0;e<M.length;e+=1){let t=ke(a,M,e),u=pe(t);de.set(u,b[e]=De(u,t))}return{c(){l=c("h3"),s=y("Delete ("),h=y(o),r=y(")"),n=k(),d=c("div"),D=c("p"),$=y("Delete a single "),F=c("strong"),G=y(q),se=y(" record."),K=k(),$e(C.$$.fragment),N=k(),P=c("h6"),P.textContent="API details",V=k(),g=c("div"),L=c("strong"),L.textContent="DELETE",ae=k(),S=c("div"),E=c("p"),oe=y("/api/collections/"),W=c("strong"),X=y(U),ne=y("/records/"),Y=c("strong"),Y.textContent=":id",ie=k(),_&&_.c(),Z=k(),T=c("div"),T.textContent="Path parameters",J=k(),B=c("table"),B.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to delete.</td></tr></tbody>',Q=k(),I=c("div"),I.textContent="Responses",x=k(),w=c("div"),O=c("div");for(let e=0;e<v.length;e+=1)v[e].c();re=k(),A=c("div");for(let e=0;e<b.length;e+=1)b[e].c();m(l,"class","m-b-sm"),m(d,"class","content txt-lg m-b-sm"),m(P,"class","m-b-xs"),m(L,"class","label label-primary"),m(S,"class","content"),m(g,"class","alert alert-danger"),m(T,"class","section-title"),m(B,"class","table-compact table-border m-b-base"),m(I,"class","section-title"),m(O,"class","tabs-header compact combined left"),m(A,"class","tabs-content"),m(w,"class","tabs")},m(e,t){p(e,l,t),i(l,s),i(l,h),i(l,r),p(e,n,t),p(e,d,t),i(d,D),i(D,$),i(D,F),i(F,G),i(D,se),p(e,K,t),Ce(C,e,t),p(e,N,t),p(e,P,t),p(e,V,t),p(e,g,t),i(g,L),i(g,ae),i(g,S),i(S,E),i(E,oe),i(E,W),i(W,X),i(E,ne),i(E,Y),i(g,ie),_&&_.m(g,null),p(e,Z,t),p(e,T,t),p(e,J,t),p(e,B,t),p(e,Q,t),p(e,I,t),p(e,x,t),p(e,w,t),i(w,O);for(let u=0;u<v.length;u+=1)v[u]&&v[u].m(O,null);i(w,re),i(w,A);for(let u=0;u<b.length;u+=1)b[u]&&b[u].m(A,null);R=!0},p(e,[t]){var be,_e;(!R||t&1)&&o!==(o=e[0].name+"")&&ee(h,o),(!R||t&1)&&q!==(q=e[0].name+"")&&ee(G,q);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').delete('RECORD_ID');
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').delete('RECORD_ID');
    `),C.$set(u),(!R||t&1)&&U!==(U=e[0].name+"")&&ee(X,U),e[1]?_||(_=ve(),_.c(),_.m(g,null)):_&&(_.d(1),_=null),t&20&&(j=H(e[4]),v=he(v,t,ue,1,e,j,ce,O,Be,ye,null,ge)),t&20&&(M=H(e[4]),Ie(),b=he(b,t,pe,1,e,M,de,A,Oe,De,null,ke),Ae())},i(e){if(!R){te(C.$$.fragment,e);for(let t=0;t<M.length;t+=1)te(b[t]);R=!0}},o(e){le(C.$$.fragment,e);for(let t=0;t<b.length;t+=1)le(b[t]);R=!1},d(e){e&&(f(l),f(n),f(d),f(K),f(N),f(P),f(V),f(g),f(Z),f(T),f(J),f(B),f(Q),f(I),f(x),f(w)),we(C,e),_&&_.d();for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<b.length;t+=1)b[t].d()}}}function je(a,l,s){let o,h,{collection:r}=l,n=204,d=[];const D=$=>s(2,n=$.code);return a.$$set=$=>{"collection"in $&&s(0,r=$.collection)},a.$$.update=()=>{a.$$.dirty&1&&s(1,o=(r==null?void 0:r.deleteRule)===null),a.$$.dirty&3&&r!=null&&r.id&&(d.push({code:204,body:`
                null
            `}),d.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `}),o&&d.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only superusers can access this action.",
                      "data": {}
                    }
                `}),d.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},s(3,h=Me.getApiExampleUrl(qe.baseURL)),[r,o,n,h,d,D]}class ze extends Re{constructor(l){super(),Ee(this,l,je,Ue,Pe,{collection:0})}}export{ze as default};
