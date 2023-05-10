import{S as ke,i as be,s as ge,e as r,w as g,b as w,c as _e,f as k,g as h,h as n,m as me,x as H,N as re,O as we,k as ve,P as Ce,n as Pe,t as L,a as Y,o as _,d as pe,T as Me,C as Se,p as $e,r as Q,u as je,M as Ae}from"./index-48bae5f5.js";import{S as Be}from"./SdkTabs-db44bf92.js";function ue(a,l,o){const s=a.slice();return s[5]=l[o],s}function de(a,l,o){const s=a.slice();return s[5]=l[o],s}function fe(a,l){let o,s=l[5].code+"",m,f,i,u;function d(){return l[4](l[5])}return{key:a,first:null,c(){o=r("button"),m=g(s),f=w(),k(o,"class","tab-item"),Q(o,"active",l[1]===l[5].code),this.first=o},m(v,C){h(v,o,C),n(o,m),n(o,f),i||(u=je(o,"click",d),i=!0)},p(v,C){l=v,C&4&&s!==(s=l[5].code+"")&&H(m,s),C&6&&Q(o,"active",l[1]===l[5].code)},d(v){v&&_(o),i=!1,u()}}}function he(a,l){let o,s,m,f;return s=new Ae({props:{content:l[5].body}}),{key:a,first:null,c(){o=r("div"),_e(s.$$.fragment),m=w(),k(o,"class","tab-item"),Q(o,"active",l[1]===l[5].code),this.first=o},m(i,u){h(i,o,u),me(s,o,null),n(o,m),f=!0},p(i,u){l=i;const d={};u&4&&(d.content=l[5].body),s.$set(d),(!f||u&6)&&Q(o,"active",l[1]===l[5].code)},i(i){f||(L(s.$$.fragment,i),f=!0)},o(i){Y(s.$$.fragment,i),f=!1},d(i){i&&_(o),pe(s)}}}function Oe(a){var ae,ne;let l,o,s=a[0].name+"",m,f,i,u,d,v,C,F=a[0].name+"",U,X,q,P,D,j,W,M,K,R,y,A,Z,V,z=a[0].name+"",E,x,I,B,J,S,O,b=[],ee=new Map,te,T,p=[],le=new Map,$;P=new Be({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        const result = await pb.collection('${(ae=a[0])==null?void 0:ae.name}').listAuthMethods();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        final result = await pb.collection('${(ne=a[0])==null?void 0:ne.name}').listAuthMethods();
    `}});let G=a[2];const oe=e=>e[5].code;for(let e=0;e<G.length;e+=1){let t=de(a,G,e),c=oe(t);ee.set(c,b[e]=fe(c,t))}let N=a[2];const se=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=ue(a,N,e),c=se(t);le.set(c,p[e]=he(c,t))}return{c(){l=r("h3"),o=g("List auth methods ("),m=g(s),f=g(")"),i=w(),u=r("div"),d=r("p"),v=g("Returns a public list with all allowed "),C=r("strong"),U=g(F),X=g(" authentication methods."),q=w(),_e(P.$$.fragment),D=w(),j=r("h6"),j.textContent="API details",W=w(),M=r("div"),K=r("strong"),K.textContent="GET",R=w(),y=r("div"),A=r("p"),Z=g("/api/collections/"),V=r("strong"),E=g(z),x=g("/auth-methods"),I=w(),B=r("div"),B.textContent="Responses",J=w(),S=r("div"),O=r("div");for(let e=0;e<b.length;e+=1)b[e].c();te=w(),T=r("div");for(let e=0;e<p.length;e+=1)p[e].c();k(l,"class","m-b-sm"),k(u,"class","content txt-lg m-b-sm"),k(j,"class","m-b-xs"),k(K,"class","label label-primary"),k(y,"class","content"),k(M,"class","alert alert-info"),k(B,"class","section-title"),k(O,"class","tabs-header compact left"),k(T,"class","tabs-content"),k(S,"class","tabs")},m(e,t){h(e,l,t),n(l,o),n(l,m),n(l,f),h(e,i,t),h(e,u,t),n(u,d),n(d,v),n(d,C),n(C,U),n(d,X),h(e,q,t),me(P,e,t),h(e,D,t),h(e,j,t),h(e,W,t),h(e,M,t),n(M,K),n(M,R),n(M,y),n(y,A),n(A,Z),n(A,V),n(V,E),n(A,x),h(e,I,t),h(e,B,t),h(e,J,t),h(e,S,t),n(S,O);for(let c=0;c<b.length;c+=1)b[c]&&b[c].m(O,null);n(S,te),n(S,T);for(let c=0;c<p.length;c+=1)p[c]&&p[c].m(T,null);$=!0},p(e,[t]){var ie,ce;(!$||t&1)&&s!==(s=e[0].name+"")&&H(m,s),(!$||t&1)&&F!==(F=e[0].name+"")&&H(U,F);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const result = await pb.collection('${(ie=e[0])==null?void 0:ie.name}').listAuthMethods();
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final result = await pb.collection('${(ce=e[0])==null?void 0:ce.name}').listAuthMethods();
    `),P.$set(c),(!$||t&1)&&z!==(z=e[0].name+"")&&H(E,z),t&6&&(G=e[2],b=re(b,t,oe,1,e,G,ee,O,we,fe,null,de)),t&6&&(N=e[2],ve(),p=re(p,t,se,1,e,N,le,T,Ce,he,null,ue),Pe())},i(e){if(!$){L(P.$$.fragment,e);for(let t=0;t<N.length;t+=1)L(p[t]);$=!0}},o(e){Y(P.$$.fragment,e);for(let t=0;t<p.length;t+=1)Y(p[t]);$=!1},d(e){e&&_(l),e&&_(i),e&&_(u),e&&_(q),pe(P,e),e&&_(D),e&&_(j),e&&_(W),e&&_(M),e&&_(I),e&&_(B),e&&_(J),e&&_(S);for(let t=0;t<b.length;t+=1)b[t].d();for(let t=0;t<p.length;t+=1)p[t].d()}}}function Te(a,l,o){let s,{collection:m=new Me}=l,f=200,i=[];const u=d=>o(1,f=d.code);return a.$$set=d=>{"collection"in d&&o(0,m=d.collection)},o(3,s=Se.getApiExampleUrl($e.baseUrl)),o(2,i=[{code:200,body:`
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
            `}]),[m,f,i,s,u]}class Ke extends ke{constructor(l){super(),be(this,l,Te,Oe,ge,{collection:0})}}export{Ke as default};
