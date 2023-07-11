import{S as Ee,i as Be,s as Se,N,e as r,w as v,b as k,c as Ce,f as b,g as d,h as n,m as ye,x as O,P as ve,Q as Me,k as Re,R as Ae,n as Ue,t as ee,a as te,o as m,d as Te,U as We,C as ze,p as He,r as F,u as je,M as De}from"./index-a084d9d7.js";import{S as Le}from"./SdkTabs-ba0ec979.js";function we(o,l,a){const s=o.slice();return s[5]=l[a],s}function $e(o,l,a){const s=o.slice();return s[5]=l[a],s}function qe(o,l){let a,s=l[5].code+"",h,f,i,p;function u(){return l[4](l[5])}return{key:o,first:null,c(){a=r("button"),h=v(s),f=k(),b(a,"class","tab-item"),F(a,"active",l[1]===l[5].code),this.first=a},m($,q){d($,a,q),n(a,h),n(a,f),i||(p=je(a,"click",u),i=!0)},p($,q){l=$,q&4&&s!==(s=l[5].code+"")&&O(h,s),q&6&&F(a,"active",l[1]===l[5].code)},d($){$&&m(a),i=!1,p()}}}function Pe(o,l){let a,s,h,f;return s=new De({props:{content:l[5].body}}),{key:o,first:null,c(){a=r("div"),Ce(s.$$.fragment),h=k(),b(a,"class","tab-item"),F(a,"active",l[1]===l[5].code),this.first=a},m(i,p){d(i,a,p),ye(s,a,null),n(a,h),f=!0},p(i,p){l=i;const u={};p&4&&(u.content=l[5].body),s.$set(u),(!f||p&6)&&F(a,"active",l[1]===l[5].code)},i(i){f||(ee(s.$$.fragment,i),f=!0)},o(i){te(s.$$.fragment,i),f=!1},d(i){i&&m(a),Te(s)}}}function Ne(o){var pe,ue,be,fe;let l,a,s=o[0].name+"",h,f,i,p,u,$,q,W=o[0].name+"",I,le,K,P,Q,T,G,w,z,ae,H,E,se,J,j=o[0].name+"",V,oe,ne,D,X,B,Y,S,Z,M,x,C,R,g=[],ie=new Map,ce,A,_=[],re=new Map,y;P=new Le({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').requestEmailChange('new@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(be=o[0])==null?void 0:be.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').requestEmailChange('new@example.com');
    `}});let L=N(o[2]);const de=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=$e(o,L,e),c=de(t);ie.set(c,g[e]=qe(c,t))}let U=N(o[2]);const me=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=we(o,U,e),c=me(t);re.set(c,_[e]=Pe(c,t))}return{c(){l=r("h3"),a=v("Request email change ("),h=v(s),f=v(")"),i=k(),p=r("div"),u=r("p"),$=v("Sends "),q=r("strong"),I=v(W),le=v(" email change request."),K=k(),Ce(P.$$.fragment),Q=k(),T=r("h6"),T.textContent="API details",G=k(),w=r("div"),z=r("strong"),z.textContent="POST",ae=k(),H=r("div"),E=r("p"),se=v("/api/collections/"),J=r("strong"),V=v(j),oe=v("/request-email-change"),ne=k(),D=r("p"),D.innerHTML="Requires record <code>Authorization:TOKEN</code> header",X=k(),B=r("div"),B.textContent="Body Parameters",Y=k(),S=r("table"),S.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>newEmail</span></div></td> <td><span class="label">String</span></td> <td>The new email address to send the change email request.</td></tr></tbody>',Z=k(),M=r("div"),M.textContent="Responses",x=k(),C=r("div"),R=r("div");for(let e=0;e<g.length;e+=1)g[e].c();ce=k(),A=r("div");for(let e=0;e<_.length;e+=1)_[e].c();b(l,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(T,"class","m-b-xs"),b(z,"class","label label-primary"),b(H,"class","content"),b(D,"class","txt-hint txt-sm txt-right"),b(w,"class","alert alert-success"),b(B,"class","section-title"),b(S,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(R,"class","tabs-header compact left"),b(A,"class","tabs-content"),b(C,"class","tabs")},m(e,t){d(e,l,t),n(l,a),n(l,h),n(l,f),d(e,i,t),d(e,p,t),n(p,u),n(u,$),n(u,q),n(q,I),n(u,le),d(e,K,t),ye(P,e,t),d(e,Q,t),d(e,T,t),d(e,G,t),d(e,w,t),n(w,z),n(w,ae),n(w,H),n(H,E),n(E,se),n(E,J),n(J,V),n(E,oe),n(w,ne),n(w,D),d(e,X,t),d(e,B,t),d(e,Y,t),d(e,S,t),d(e,Z,t),d(e,M,t),d(e,x,t),d(e,C,t),n(C,R);for(let c=0;c<g.length;c+=1)g[c]&&g[c].m(R,null);n(C,ce),n(C,A);for(let c=0;c<_.length;c+=1)_[c]&&_[c].m(A,null);y=!0},p(e,[t]){var he,_e,ke,ge;(!y||t&1)&&s!==(s=e[0].name+"")&&O(h,s),(!y||t&1)&&W!==(W=e[0].name+"")&&O(I,W);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(he=e[0])==null?void 0:he.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').requestEmailChange('new@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(ke=e[0])==null?void 0:ke.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(ge=e[0])==null?void 0:ge.name}').requestEmailChange('new@example.com');
    `),P.$set(c),(!y||t&1)&&j!==(j=e[0].name+"")&&O(V,j),t&6&&(L=N(e[2]),g=ve(g,t,de,1,e,L,ie,R,Me,qe,null,$e)),t&6&&(U=N(e[2]),Re(),_=ve(_,t,me,1,e,U,re,A,Ae,Pe,null,we),Ue())},i(e){if(!y){ee(P.$$.fragment,e);for(let t=0;t<U.length;t+=1)ee(_[t]);y=!0}},o(e){te(P.$$.fragment,e);for(let t=0;t<_.length;t+=1)te(_[t]);y=!1},d(e){e&&(m(l),m(i),m(p),m(K),m(Q),m(T),m(G),m(w),m(X),m(B),m(Y),m(S),m(Z),m(M),m(x),m(C)),Te(P,e);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function Oe(o,l,a){let s,{collection:h=new We}=l,f=204,i=[];const p=u=>a(1,f=u.code);return o.$$set=u=>{"collection"in u&&a(0,h=u.collection)},a(3,s=ze.getApiExampleUrl(He.baseUrl)),a(2,i=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "newEmail": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `},{code:401,body:`
                {
                  "code": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `},{code:403,body:`
                {
                  "code": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `}]),[h,f,i,s,p]}class Ke extends Ee{constructor(l){super(),Be(this,l,Oe,Ne,Se,{collection:0})}}export{Ke as default};
