import{S as Re,i as Pe,s as Ee,N as j,e as c,w as y,b as k,c as De,f as m,g as u,h as i,m as we,x as ee,P as he,Q as Te,k as Be,R as Oe,n as Ie,t as te,a as le,o as p,d as Ce,U as Me,C as Ae,p as Se,r as N,u as qe,M as Ue}from"./index-a084d9d7.js";import{S as He}from"./SdkTabs-ba0ec979.js";function ke(a,l,s){const n=a.slice();return n[6]=l[s],n}function ge(a,l,s){const n=a.slice();return n[6]=l[s],n}function ve(a){let l;return{c(){l=c("p"),l.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",m(l,"class","txt-hint txt-sm txt-right")},m(s,n){u(s,l,n)},d(s){s&&p(l)}}}function ye(a,l){let s,n,h;function d(){return l[5](l[6])}return{key:a,first:null,c(){s=c("button"),s.textContent=`${l[6].code} `,m(s,"class","tab-item"),N(s,"active",l[2]===l[6].code),this.first=s},m(o,r){u(o,s,r),n||(h=qe(s,"click",d),n=!0)},p(o,r){l=o,r&20&&N(s,"active",l[2]===l[6].code)},d(o){o&&p(s),n=!1,h()}}}function $e(a,l){let s,n,h,d;return n=new Ue({props:{content:l[6].body}}),{key:a,first:null,c(){s=c("div"),De(n.$$.fragment),h=k(),m(s,"class","tab-item"),N(s,"active",l[2]===l[6].code),this.first=s},m(o,r){u(o,s,r),we(n,s,null),i(s,h),d=!0},p(o,r){l=o,(!d||r&20)&&N(s,"active",l[2]===l[6].code)},i(o){d||(te(n.$$.fragment,o),d=!0)},o(o){le(n.$$.fragment,o),d=!1},d(o){o&&p(s),Ce(n)}}}function Le(a){var pe,me;let l,s,n=a[0].name+"",h,d,o,r,$,D,z,S=a[0].name+"",F,se,K,w,Q,E,G,g,q,ae,U,P,ne,J,H=a[0].name+"",V,oe,W,ie,X,T,Y,B,Z,O,x,C,I,v=[],ce=new Map,de,M,b=[],re=new Map,R;w=new He({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(pe=a[0])==null?void 0:pe.name}').delete('RECORD_ID');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(me=a[0])==null?void 0:me.name}').delete('RECORD_ID');
    `}});let _=a[1]&&ve(),L=j(a[4]);const fe=e=>e[6].code;for(let e=0;e<L.length;e+=1){let t=ge(a,L,e),f=fe(t);ce.set(f,v[e]=ye(f,t))}let A=j(a[4]);const ue=e=>e[6].code;for(let e=0;e<A.length;e+=1){let t=ke(a,A,e),f=ue(t);re.set(f,b[e]=$e(f,t))}return{c(){l=c("h3"),s=y("Delete ("),h=y(n),d=y(")"),o=k(),r=c("div"),$=c("p"),D=y("Delete a single "),z=c("strong"),F=y(S),se=y(" record."),K=k(),De(w.$$.fragment),Q=k(),E=c("h6"),E.textContent="API details",G=k(),g=c("div"),q=c("strong"),q.textContent="DELETE",ae=k(),U=c("div"),P=c("p"),ne=y("/api/collections/"),J=c("strong"),V=y(H),oe=y("/records/"),W=c("strong"),W.textContent=":id",ie=k(),_&&_.c(),X=k(),T=c("div"),T.textContent="Path parameters",Y=k(),B=c("table"),B.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the record to delete.</td></tr></tbody>',Z=k(),O=c("div"),O.textContent="Responses",x=k(),C=c("div"),I=c("div");for(let e=0;e<v.length;e+=1)v[e].c();de=k(),M=c("div");for(let e=0;e<b.length;e+=1)b[e].c();m(l,"class","m-b-sm"),m(r,"class","content txt-lg m-b-sm"),m(E,"class","m-b-xs"),m(q,"class","label label-primary"),m(U,"class","content"),m(g,"class","alert alert-danger"),m(T,"class","section-title"),m(B,"class","table-compact table-border m-b-base"),m(O,"class","section-title"),m(I,"class","tabs-header compact left"),m(M,"class","tabs-content"),m(C,"class","tabs")},m(e,t){u(e,l,t),i(l,s),i(l,h),i(l,d),u(e,o,t),u(e,r,t),i(r,$),i($,D),i($,z),i(z,F),i($,se),u(e,K,t),we(w,e,t),u(e,Q,t),u(e,E,t),u(e,G,t),u(e,g,t),i(g,q),i(g,ae),i(g,U),i(U,P),i(P,ne),i(P,J),i(J,V),i(P,oe),i(P,W),i(g,ie),_&&_.m(g,null),u(e,X,t),u(e,T,t),u(e,Y,t),u(e,B,t),u(e,Z,t),u(e,O,t),u(e,x,t),u(e,C,t),i(C,I);for(let f=0;f<v.length;f+=1)v[f]&&v[f].m(I,null);i(C,de),i(C,M);for(let f=0;f<b.length;f+=1)b[f]&&b[f].m(M,null);R=!0},p(e,[t]){var be,_e;(!R||t&1)&&n!==(n=e[0].name+"")&&ee(h,n),(!R||t&1)&&S!==(S=e[0].name+"")&&ee(F,S);const f={};t&9&&(f.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').delete('RECORD_ID');
    `),t&9&&(f.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').delete('RECORD_ID');
    `),w.$set(f),(!R||t&1)&&H!==(H=e[0].name+"")&&ee(V,H),e[1]?_||(_=ve(),_.c(),_.m(g,null)):_&&(_.d(1),_=null),t&20&&(L=j(e[4]),v=he(v,t,fe,1,e,L,ce,I,Te,ye,null,ge)),t&20&&(A=j(e[4]),Be(),b=he(b,t,ue,1,e,A,re,M,Oe,$e,null,ke),Ie())},i(e){if(!R){te(w.$$.fragment,e);for(let t=0;t<A.length;t+=1)te(b[t]);R=!0}},o(e){le(w.$$.fragment,e);for(let t=0;t<b.length;t+=1)le(b[t]);R=!1},d(e){e&&(p(l),p(o),p(r),p(K),p(Q),p(E),p(G),p(g),p(X),p(T),p(Y),p(B),p(Z),p(O),p(x),p(C)),Ce(w,e),_&&_.d();for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<b.length;t+=1)b[t].d()}}}function je(a,l,s){let n,h,{collection:d=new Me}=l,o=204,r=[];const $=D=>s(2,o=D.code);return a.$$set=D=>{"collection"in D&&s(0,d=D.collection)},a.$$.update=()=>{a.$$.dirty&1&&s(1,n=(d==null?void 0:d.deleteRule)===null),a.$$.dirty&3&&d!=null&&d.id&&(r.push({code:204,body:`
                null
            `}),r.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `}),n&&r.push({code:403,body:`
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
            `}))},s(3,h=Ae.getApiExampleUrl(Se.baseUrl)),[d,n,o,h,r,$]}class Fe extends Re{constructor(l){super(),Pe(this,l,je,Le,Ee,{collection:0})}}export{Fe as default};
