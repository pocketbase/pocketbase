import{S as Pe,i as Se,s as Te,e as r,w,b as k,c as ge,f as b,g as f,h as i,m as ye,x as D,aa as _e,ab as Be,k as qe,ac as Re,n as Ee,t as Z,a as x,o as p,d as Ce,ae as Oe,C as Ve,p as Ke,r as H,u as Me,a9 as Ne}from"./index-76001414.js";import{S as Ae}from"./SdkTabs-9688a076.js";function ke(a,l,s){const o=a.slice();return o[5]=l[s],o}function ve(a,l,s){const o=a.slice();return o[5]=l[s],o}function we(a,l){let s,o=l[5].code+"",h,d,n,u;function m(){return l[4](l[5])}return{key:a,first:null,c(){s=r("button"),h=w(o),d=k(),b(s,"class","tab-item"),H(s,"active",l[1]===l[5].code),this.first=s},m($,g){f($,s,g),i(s,h),i(s,d),n||(u=Me(s,"click",m),n=!0)},p($,g){l=$,g&4&&o!==(o=l[5].code+"")&&D(h,o),g&6&&H(s,"active",l[1]===l[5].code)},d($){$&&p(s),n=!1,u()}}}function $e(a,l){let s,o,h,d;return o=new Ne({props:{content:l[5].body}}),{key:a,first:null,c(){s=r("div"),ge(o.$$.fragment),h=k(),b(s,"class","tab-item"),H(s,"active",l[1]===l[5].code),this.first=s},m(n,u){f(n,s,u),ye(o,s,null),i(s,h),d=!0},p(n,u){l=n;const m={};u&4&&(m.content=l[5].body),o.$set(m),(!d||u&6)&&H(s,"active",l[1]===l[5].code)},i(n){d||(Z(o.$$.fragment,n),d=!0)},o(n){x(o.$$.fragment,n),d=!1},d(n){n&&p(s),Ce(o)}}}function Ue(a){var re,fe,pe,ue;let l,s,o=a[0].name+"",h,d,n,u,m,$,g,M=a[0].name+"",F,ee,I,y,L,T,z,C,N,te,A,B,le,G,U=a[0].name+"",J,se,Q,q,W,R,X,E,Y,P,O,v=[],oe=new Map,ae,V,_=[],ie=new Map,S;y=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(re=a[0])==null?void 0:re.name}').confirmVerification('TOKEN');

        // optionally refresh the previous authStore state with the latest record changes
        await pb.collection('${(fe=a[0])==null?void 0:fe.name}').authRefresh();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(pe=a[0])==null?void 0:pe.name}').confirmVerification('TOKEN');

        // optionally refresh the previous authStore state with the latest record changes
        await pb.collection('${(ue=a[0])==null?void 0:ue.name}').authRefresh();
    `}});let j=a[2];const ne=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=ve(a,j,e),c=ne(t);oe.set(c,v[e]=we(c,t))}let K=a[2];const ce=e=>e[5].code;for(let e=0;e<K.length;e+=1){let t=ke(a,K,e),c=ce(t);ie.set(c,_[e]=$e(c,t))}return{c(){l=r("h3"),s=w("Confirm verification ("),h=w(o),d=w(")"),n=k(),u=r("div"),m=r("p"),$=w("Confirms "),g=r("strong"),F=w(M),ee=w(" account verification request."),I=k(),ge(y.$$.fragment),L=k(),T=r("h6"),T.textContent="API details",z=k(),C=r("div"),N=r("strong"),N.textContent="POST",te=k(),A=r("div"),B=r("p"),le=w("/api/collections/"),G=r("strong"),J=w(U),se=w("/confirm-verification"),Q=k(),q=r("div"),q.textContent="Body Parameters",W=k(),R=r("table"),R.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the verification request email.</td></tr></tbody>`,X=k(),E=r("div"),E.textContent="Responses",Y=k(),P=r("div"),O=r("div");for(let e=0;e<v.length;e+=1)v[e].c();ae=k(),V=r("div");for(let e=0;e<_.length;e+=1)_[e].c();b(l,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(T,"class","m-b-xs"),b(N,"class","label label-primary"),b(A,"class","content"),b(C,"class","alert alert-success"),b(q,"class","section-title"),b(R,"class","table-compact table-border m-b-base"),b(E,"class","section-title"),b(O,"class","tabs-header compact left"),b(V,"class","tabs-content"),b(P,"class","tabs")},m(e,t){f(e,l,t),i(l,s),i(l,h),i(l,d),f(e,n,t),f(e,u,t),i(u,m),i(m,$),i(m,g),i(g,F),i(m,ee),f(e,I,t),ye(y,e,t),f(e,L,t),f(e,T,t),f(e,z,t),f(e,C,t),i(C,N),i(C,te),i(C,A),i(A,B),i(B,le),i(B,G),i(G,J),i(B,se),f(e,Q,t),f(e,q,t),f(e,W,t),f(e,R,t),f(e,X,t),f(e,E,t),f(e,Y,t),f(e,P,t),i(P,O);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(O,null);i(P,ae),i(P,V);for(let c=0;c<_.length;c+=1)_[c]&&_[c].m(V,null);S=!0},p(e,[t]){var me,de,be,he;(!S||t&1)&&o!==(o=e[0].name+"")&&D(h,o),(!S||t&1)&&M!==(M=e[0].name+"")&&D(F,M);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').confirmVerification('TOKEN');

        // optionally refresh the previous authStore state with the latest record changes
        await pb.collection('${(de=e[0])==null?void 0:de.name}').authRefresh();
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').confirmVerification('TOKEN');

        // optionally refresh the previous authStore state with the latest record changes
        await pb.collection('${(he=e[0])==null?void 0:he.name}').authRefresh();
    `),y.$set(c),(!S||t&1)&&U!==(U=e[0].name+"")&&D(J,U),t&6&&(j=e[2],v=_e(v,t,ne,1,e,j,oe,O,Be,we,null,ve)),t&6&&(K=e[2],qe(),_=_e(_,t,ce,1,e,K,ie,V,Re,$e,null,ke),Ee())},i(e){if(!S){Z(y.$$.fragment,e);for(let t=0;t<K.length;t+=1)Z(_[t]);S=!0}},o(e){x(y.$$.fragment,e);for(let t=0;t<_.length;t+=1)x(_[t]);S=!1},d(e){e&&p(l),e&&p(n),e&&p(u),e&&p(I),Ce(y,e),e&&p(L),e&&p(T),e&&p(z),e&&p(C),e&&p(Q),e&&p(q),e&&p(W),e&&p(R),e&&p(X),e&&p(E),e&&p(Y),e&&p(P);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function je(a,l,s){let o,{collection:h=new Oe}=l,d=204,n=[];const u=m=>s(1,d=m.code);return a.$$set=m=>{"collection"in m&&s(0,h=m.collection)},s(3,o=Ve.getApiExampleUrl(Ke.baseUrl)),s(2,n=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "token": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}]),[h,d,n,o,u]}class Fe extends Pe{constructor(l){super(),Se(this,l,je,Ue,Te,{collection:0})}}export{Fe as default};
