import{S as $e,i as Se,s as ye,e as c,w,b as k,c as oe,f as h,g as d,h as a,m as se,x as G,N as we,P as Te,k as je,Q as Ae,n as Be,t as U,a as W,o as u,d as ae,T as Fe,C as Oe,p as Qe,r as V,u as Ne,M as He}from"./index-a65ca895.js";import{S as Ke}from"./SdkTabs-ad912c8f.js";import{F as qe}from"./FieldsQueryParam-ba250473.js";function ve(n,l,o){const s=n.slice();return s[5]=l[o],s}function Ce(n,l,o){const s=n.slice();return s[5]=l[o],s}function Pe(n,l){let o,s=l[5].code+"",_,p,i,f;function m(){return l[4](l[5])}return{key:n,first:null,c(){o=c("button"),_=w(s),p=k(),h(o,"class","tab-item"),V(o,"active",l[1]===l[5].code),this.first=o},m(v,C){d(v,o,C),a(o,_),a(o,p),i||(f=Ne(o,"click",m),i=!0)},p(v,C){l=v,C&4&&s!==(s=l[5].code+"")&&G(_,s),C&6&&V(o,"active",l[1]===l[5].code)},d(v){v&&u(o),i=!1,f()}}}function Me(n,l){let o,s,_,p;return s=new He({props:{content:l[5].body}}),{key:n,first:null,c(){o=c("div"),oe(s.$$.fragment),_=k(),h(o,"class","tab-item"),V(o,"active",l[1]===l[5].code),this.first=o},m(i,f){d(i,o,f),se(s,o,null),a(o,_),p=!0},p(i,f){l=i;const m={};f&4&&(m.content=l[5].body),s.$set(m),(!p||f&6)&&V(o,"active",l[1]===l[5].code)},i(i){p||(U(s.$$.fragment,i),p=!0)},o(i){W(s.$$.fragment,i),p=!1},d(i){i&&u(o),ae(s)}}}function ze(n){var _e,be;let l,o,s=n[0].name+"",_,p,i,f,m,v,C,H=n[0].name+"",L,ne,E,P,I,j,J,M,K,ie,q,A,ce,Y,z=n[0].name+"",X,re,R,B,Z,$,x,de,ee,T,te,F,le,S,O,g=[],ue=new Map,fe,Q,b=[],me=new Map,y;P=new Ke({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        const result = await pb.collection('${(_e=n[0])==null?void 0:_e.name}').listAuthMethods();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        final result = await pb.collection('${(be=n[0])==null?void 0:be.name}').listAuthMethods();
    `}}),T=new qe({});let D=n[2];const pe=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=Ce(n,D,e),r=pe(t);ue.set(r,g[e]=Pe(r,t))}let N=n[2];const he=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=ve(n,N,e),r=he(t);me.set(r,b[e]=Me(r,t))}return{c(){l=c("h3"),o=w("List auth methods ("),_=w(s),p=w(")"),i=k(),f=c("div"),m=c("p"),v=w("Returns a public list with all allowed "),C=c("strong"),L=w(H),ne=w(" authentication methods."),E=k(),oe(P.$$.fragment),I=k(),j=c("h6"),j.textContent="API details",J=k(),M=c("div"),K=c("strong"),K.textContent="GET",ie=k(),q=c("div"),A=c("p"),ce=w("/api/collections/"),Y=c("strong"),X=w(z),re=w("/auth-methods"),R=k(),B=c("div"),B.textContent="Query parameters",Z=k(),$=c("table"),x=c("thead"),x.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,de=k(),ee=c("tbody"),oe(T.$$.fragment),te=k(),F=c("div"),F.textContent="Responses",le=k(),S=c("div"),O=c("div");for(let e=0;e<g.length;e+=1)g[e].c();fe=k(),Q=c("div");for(let e=0;e<b.length;e+=1)b[e].c();h(l,"class","m-b-sm"),h(f,"class","content txt-lg m-b-sm"),h(j,"class","m-b-xs"),h(K,"class","label label-primary"),h(q,"class","content"),h(M,"class","alert alert-info"),h(B,"class","section-title"),h($,"class","table-compact table-border m-b-base"),h(F,"class","section-title"),h(O,"class","tabs-header compact left"),h(Q,"class","tabs-content"),h(S,"class","tabs")},m(e,t){d(e,l,t),a(l,o),a(l,_),a(l,p),d(e,i,t),d(e,f,t),a(f,m),a(m,v),a(m,C),a(C,L),a(m,ne),d(e,E,t),se(P,e,t),d(e,I,t),d(e,j,t),d(e,J,t),d(e,M,t),a(M,K),a(M,ie),a(M,q),a(q,A),a(A,ce),a(A,Y),a(Y,X),a(A,re),d(e,R,t),d(e,B,t),d(e,Z,t),d(e,$,t),a($,x),a($,de),a($,ee),se(T,ee,null),d(e,te,t),d(e,F,t),d(e,le,t),d(e,S,t),a(S,O);for(let r=0;r<g.length;r+=1)g[r]&&g[r].m(O,null);a(S,fe),a(S,Q);for(let r=0;r<b.length;r+=1)b[r]&&b[r].m(Q,null);y=!0},p(e,[t]){var ke,ge;(!y||t&1)&&s!==(s=e[0].name+"")&&G(_,s),(!y||t&1)&&H!==(H=e[0].name+"")&&G(L,H);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const result = await pb.collection('${(ke=e[0])==null?void 0:ke.name}').listAuthMethods();
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final result = await pb.collection('${(ge=e[0])==null?void 0:ge.name}').listAuthMethods();
    `),P.$set(r),(!y||t&1)&&z!==(z=e[0].name+"")&&G(X,z),t&6&&(D=e[2],g=we(g,t,pe,1,e,D,ue,O,Te,Pe,null,Ce)),t&6&&(N=e[2],je(),b=we(b,t,he,1,e,N,me,Q,Ae,Me,null,ve),Be())},i(e){if(!y){U(P.$$.fragment,e),U(T.$$.fragment,e);for(let t=0;t<N.length;t+=1)U(b[t]);y=!0}},o(e){W(P.$$.fragment,e),W(T.$$.fragment,e);for(let t=0;t<b.length;t+=1)W(b[t]);y=!1},d(e){e&&u(l),e&&u(i),e&&u(f),e&&u(E),ae(P,e),e&&u(I),e&&u(j),e&&u(J),e&&u(M),e&&u(R),e&&u(B),e&&u(Z),e&&u($),ae(T),e&&u(te),e&&u(F),e&&u(le),e&&u(S);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<b.length;t+=1)b[t].d()}}}function De(n,l,o){let s,{collection:_=new Fe}=l,p=200,i=[];const f=m=>o(1,p=m.code);return n.$$set=m=>{"collection"in m&&o(0,_=m.collection)},o(3,s=Oe.getApiExampleUrl(Qe.baseUrl)),o(2,i=[{code:200,body:`
                {
                  "usernamePassword": true,
                  "emailPassword": true,
                  "authProviders": [
                    {
                      "name": "github",
                      "state": "3Yd8jNkK_6PJG6hPWwBjLqKwse6Ejd",
                      "codeVerifier": "KxFDWz1B3fxscCDJ_9gHQhLuh__ie7",
                      "codeChallenge": "NM1oVexB6Q6QH8uPtOUfK7tq4pmu4Jz6lNDIwoxHZNE=",
                      "codeChallengeMethod": "S256",
                      "authUrl": "https://github.com/login/oauth/authorize?client_id=demo&code_challenge=NM1oVexB6Q6QH8uPtOUfK7tq4pmu4Jz6lNDIwoxHZNE%3D&code_challenge_method=S256&response_type=code&scope=user&state=3Yd8jNkK_6PJG6hPWwBjLqKwse6Ejd&redirect_uri="
                    },
                    {
                      "name": "gitlab",
                      "state": "NeQSbtO5cShr_mk5__3CUukiMnymeb",
                      "codeVerifier": "ahTFHOgua8mkvPAlIBGwCUJbWKR_xi",
                      "codeChallenge": "O-GATkTj4eXDCnfonsqGLCd6njvTixlpCMvy5kjgOOg=",
                      "codeChallengeMethod": "S256",
                      "authUrl": "https://gitlab.com/oauth/authorize?client_id=demo&code_challenge=O-GATkTj4eXDCnfonsqGLCd6njvTixlpCMvy5kjgOOg%3D&code_challenge_method=S256&response_type=code&scope=read_user&state=NeQSbtO5cShr_mk5__3CUukiMnymeb&redirect_uri="
                    },
                    {
                      "name": "google",
                      "state": "zB3ZPifV1TW2GMuvuFkamSXfSNkHPQ",
                      "codeVerifier": "t3CmO5VObGzdXqieakvR_fpjiW0zdO",
                      "codeChallenge": "KChwoQPKYlz2anAdqtgsSTdIo8hdwtc1fh2wHMwW2Yk=",
                      "codeChallengeMethod": "S256",
                      "authUrl": "https://accounts.google.com/o/oauth2/auth?client_id=demo&code_challenge=KChwoQPKYlz2anAdqtgsSTdIo8hdwtc1fh2wHMwW2Yk%3D&code_challenge_method=S256&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&state=zB3ZPifV1TW2GMuvuFkamSXfSNkHPQ&redirect_uri="
                    }
                  ]
                }
            `}]),[_,p,i,s,f]}class Ve extends $e{constructor(l){super(),Se(this,l,De,ze,ye,{collection:0})}}export{Ve as default};
