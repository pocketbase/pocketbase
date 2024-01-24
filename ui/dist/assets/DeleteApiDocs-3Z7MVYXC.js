import{S as Re,i as Pe,s as Ee,O as j,e as c,w as y,b as k,c as De,f as m,g as p,h as i,m as Ce,x as ee,P as he,Q as Oe,k as Te,R as Be,n as Ie,t as te,a as le,o as u,d as we,C as Ae,p as Me,r as N,u as Se,N as qe}from"./index-78piLIP3.js";import{S as He}from"./SdkTabs-c6VuPJvR.js";function ke(a,l,s){const o=a.slice();return o[6]=l[s],o}function ge(a,l,s){const o=a.slice();return o[6]=l[s],o}function ve(a){let l;return{c(){l=c("p"),l.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",m(l,"class","txt-hint txt-sm txt-right")},m(s,o){p(s,l,o)},d(s){s&&u(l)}}}function ye(a,l){let s,o,h;function d(){return l[5](l[6])}return{key:a,first:null,c(){s=c("button"),s.textContent=`${l[6].code} `,m(s,"class","tab-item"),N(s,"active",l[2]===l[6].code),this.first=s},m(n,r){p(n,s,r),o||(h=Se(s,"click",d),o=!0)},p(n,r){l=n,r&20&&N(s,"active",l[2]===l[6].code)},d(n){n&&u(s),o=!1,h()}}}function $e(a,l){let s,o,h,d;return o=new qe({props:{content:l[6].body}}),{key:a,first:null,c(){s=c("div"),De(o.$$.fragment),h=k(),m(s,"class","tab-item"),N(s,"active",l[2]===l[6].code),this.first=s},m(n,r){p(n,s,r),Ce(o,s,null),i(s,h),d=!0},p(n,r){l=n,(!d||r&20)&&N(s,"active",l[2]===l[6].code)},i(n){d||(te(o.$$.fragment,n),d=!0)},o(n){le(o.$$.fragment,n),d=!1},d(n){n&&u(s),we(o)}}}function Le(a){var ue,me;let l,s,o=a[0].name+"",h,d,n,r,$,D,z,S=a[0].name+"",F,se,K,C,Q,E,G,g,q,ae,H,P,oe,J,L=a[0].name+"",V,ne,W,ie,X,O,Y,T,Z,B,x,w,I,v=[],ce=new Map,de,A,b=[],re=new Map,R;C=new He({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(ue=a[0])==null?void 0:ue.name}').delete('RECORD_ID');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(me=a[0])==null?void 0:me.name}').delete('RECORD_ID');
    `}});let _=a[1]&&ve(),U=j(a[4]);const fe=e=>e[6].code;for(let e=0;e<U.length;e+=1){let t=ge(a,U,e),f=fe(t);ce.set(f,v[e]=ye(f,t))}let M=j(a[4]);const pe=e=>e[6].code;for(let e=0;e<M.length;e+=1){let t=ke(a,M,e),f=pe(t);re.set(f,b[e]=$e(f,t))}return{c(){l=c("h3"),s=y("Delete ("),h=y(o),d=y(")"),n=k(),r=c("div"),$=c("p"),D=y("Delete a single "),z=c("strong"),F=y(S),se=y(" record."),K=k(),De(C.$$.fragment),Q=k(),E=c("h6"),E.textContent="API details",G=k(),g=c("div"),q=c("strong"),q.textContent="DELETE",ae=k(),H=c("div"),P=c("p"),oe=y("/api/collections/"),J=c("strong"),V=y(L),ne=y("/records/"),W=c("strong"),W.textContent=":id",ie=k(),_&&_.c(),X=k(),O=c("div"),O.textContent="Path parameters",Y=k(),T=c("table"),T.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to delete.</td></tr></tbody>',Z=k(),B=c("div"),B.textContent="Responses",x=k(),w=c("div"),I=c("div");for(let e=0;e<v.length;e+=1)v[e].c();de=k(),A=c("div");for(let e=0;e<b.length;e+=1)b[e].c();m(l,"class","m-b-sm"),m(r,"class","content txt-lg m-b-sm"),m(E,"class","m-b-xs"),m(q,"class","label label-primary"),m(H,"class","content"),m(g,"class","alert alert-danger"),m(O,"class","section-title"),m(T,"class","table-compact table-border m-b-base"),m(B,"class","section-title"),m(I,"class","tabs-header compact combined left"),m(A,"class","tabs-content"),m(w,"class","tabs")},m(e,t){p(e,l,t),i(l,s),i(l,h),i(l,d),p(e,n,t),p(e,r,t),i(r,$),i($,D),i($,z),i(z,F),i($,se),p(e,K,t),Ce(C,e,t),p(e,Q,t),p(e,E,t),p(e,G,t),p(e,g,t),i(g,q),i(g,ae),i(g,H),i(H,P),i(P,oe),i(P,J),i(J,V),i(P,ne),i(P,W),i(g,ie),_&&_.m(g,null),p(e,X,t),p(e,O,t),p(e,Y,t),p(e,T,t),p(e,Z,t),p(e,B,t),p(e,x,t),p(e,w,t),i(w,I);for(let f=0;f<v.length;f+=1)v[f]&&v[f].m(I,null);i(w,de),i(w,A);for(let f=0;f<b.length;f+=1)b[f]&&b[f].m(A,null);R=!0},p(e,[t]){var be,_e;(!R||t&1)&&o!==(o=e[0].name+"")&&ee(h,o),(!R||t&1)&&S!==(S=e[0].name+"")&&ee(F,S);const f={};t&9&&(f.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').delete('RECORD_ID');
    `),t&9&&(f.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').delete('RECORD_ID');
    `),C.$set(f),(!R||t&1)&&L!==(L=e[0].name+"")&&ee(V,L),e[1]?_||(_=ve(),_.c(),_.m(g,null)):_&&(_.d(1),_=null),t&20&&(U=j(e[4]),v=he(v,t,fe,1,e,U,ce,I,Oe,ye,null,ge)),t&20&&(M=j(e[4]),Te(),b=he(b,t,pe,1,e,M,re,A,Be,$e,null,ke),Ie())},i(e){if(!R){te(C.$$.fragment,e);for(let t=0;t<M.length;t+=1)te(b[t]);R=!0}},o(e){le(C.$$.fragment,e);for(let t=0;t<b.length;t+=1)le(b[t]);R=!1},d(e){e&&(u(l),u(n),u(r),u(K),u(Q),u(E),u(G),u(g),u(X),u(O),u(Y),u(T),u(Z),u(B),u(x),u(w)),we(C,e),_&&_.d();for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<b.length;t+=1)b[t].d()}}}function Ue(a,l,s){let o,h,{collection:d}=l,n=204,r=[];const $=D=>s(2,n=D.code);return a.$$set=D=>{"collection"in D&&s(0,d=D.collection)},a.$$.update=()=>{a.$$.dirty&1&&s(1,o=(d==null?void 0:d.deleteRule)===null),a.$$.dirty&3&&d!=null&&d.id&&(r.push({code:204,body:`
                null
            `}),r.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `}),o&&r.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}),r.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},s(3,h=Ae.getApiExampleUrl(Me.baseUrl)),[d,o,n,h,r,$]}class ze extends Re{constructor(l){super(),Pe(this,l,Ue,Le,Ee,{collection:0})}}export{ze as default};
