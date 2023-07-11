import{S as Se,i as ye,s as Te,N as D,e as c,w,b as k,c as se,f as p,g as d,h as a,m as ae,x as G,P as ve,Q as je,k as Ae,R as Be,n as Fe,t as W,a as V,o as u,d as ne,U as Oe,C as Qe,p as Ne,r as L,u as He,M as Ke}from"./index-a084d9d7.js";import{S as Ue}from"./SdkTabs-ba0ec979.js";import{F as qe}from"./FieldsQueryParam-71e01e64.js";function Ce(n,l,o){const s=n.slice();return s[5]=l[o],s}function Pe(n,l,o){const s=n.slice();return s[5]=l[o],s}function Me(n,l){let o,s=l[5].code+"",_,f,i,h;function m(){return l[4](l[5])}return{key:n,first:null,c(){o=c("button"),_=w(s),f=k(),p(o,"class","tab-item"),L(o,"active",l[1]===l[5].code),this.first=o},m(v,C){d(v,o,C),a(o,_),a(o,f),i||(h=He(o,"click",m),i=!0)},p(v,C){l=v,C&4&&s!==(s=l[5].code+"")&&G(_,s),C&6&&L(o,"active",l[1]===l[5].code)},d(v){v&&u(o),i=!1,h()}}}function $e(n,l){let o,s,_,f;return s=new Ke({props:{content:l[5].body}}),{key:n,first:null,c(){o=c("div"),se(s.$$.fragment),_=k(),p(o,"class","tab-item"),L(o,"active",l[1]===l[5].code),this.first=o},m(i,h){d(i,o,h),ae(s,o,null),a(o,_),f=!0},p(i,h){l=i;const m={};h&4&&(m.content=l[5].body),s.$set(m),(!f||h&6)&&L(o,"active",l[1]===l[5].code)},i(i){f||(W(s.$$.fragment,i),f=!0)},o(i){V(s.$$.fragment,i),f=!1},d(i){i&&u(o),ne(s)}}}function ze(n){var be,ke;let l,o,s=n[0].name+"",_,f,i,h,m,v,C,H=n[0].name+"",E,ie,I,P,J,j,Y,M,K,ce,U,A,re,R,q=n[0].name+"",X,de,Z,B,x,$,ee,ue,te,T,le,F,oe,S,O,g=[],he=new Map,me,Q,b=[],fe=new Map,y;P=new Ue({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        const result = await pb.collection('${(be=n[0])==null?void 0:be.name}').listAuthMethods();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        final result = await pb.collection('${(ke=n[0])==null?void 0:ke.name}').listAuthMethods();
    `}}),T=new qe({});let z=D(n[2]);const pe=e=>e[5].code;for(let e=0;e<z.length;e+=1){let t=Pe(n,z,e),r=pe(t);he.set(r,g[e]=Me(r,t))}let N=D(n[2]);const _e=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=Ce(n,N,e),r=_e(t);fe.set(r,b[e]=$e(r,t))}return{c(){l=c("h3"),o=w("List auth methods ("),_=w(s),f=w(")"),i=k(),h=c("div"),m=c("p"),v=w("Returns a public list with all allowed "),C=c("strong"),E=w(H),ie=w(" authentication methods."),I=k(),se(P.$$.fragment),J=k(),j=c("h6"),j.textContent="API details",Y=k(),M=c("div"),K=c("strong"),K.textContent="GET",ce=k(),U=c("div"),A=c("p"),re=w("/api/collections/"),R=c("strong"),X=w(q),de=w("/auth-methods"),Z=k(),B=c("div"),B.textContent="Query parameters",x=k(),$=c("table"),ee=c("thead"),ee.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',ue=k(),te=c("tbody"),se(T.$$.fragment),le=k(),F=c("div"),F.textContent="Responses",oe=k(),S=c("div"),O=c("div");for(let e=0;e<g.length;e+=1)g[e].c();me=k(),Q=c("div");for(let e=0;e<b.length;e+=1)b[e].c();p(l,"class","m-b-sm"),p(h,"class","content txt-lg m-b-sm"),p(j,"class","m-b-xs"),p(K,"class","label label-primary"),p(U,"class","content"),p(M,"class","alert alert-info"),p(B,"class","section-title"),p($,"class","table-compact table-border m-b-base"),p(F,"class","section-title"),p(O,"class","tabs-header compact left"),p(Q,"class","tabs-content"),p(S,"class","tabs")},m(e,t){d(e,l,t),a(l,o),a(l,_),a(l,f),d(e,i,t),d(e,h,t),a(h,m),a(m,v),a(m,C),a(C,E),a(m,ie),d(e,I,t),ae(P,e,t),d(e,J,t),d(e,j,t),d(e,Y,t),d(e,M,t),a(M,K),a(M,ce),a(M,U),a(U,A),a(A,re),a(A,R),a(R,X),a(A,de),d(e,Z,t),d(e,B,t),d(e,x,t),d(e,$,t),a($,ee),a($,ue),a($,te),ae(T,te,null),d(e,le,t),d(e,F,t),d(e,oe,t),d(e,S,t),a(S,O);for(let r=0;r<g.length;r+=1)g[r]&&g[r].m(O,null);a(S,me),a(S,Q);for(let r=0;r<b.length;r+=1)b[r]&&b[r].m(Q,null);y=!0},p(e,[t]){var ge,we;(!y||t&1)&&s!==(s=e[0].name+"")&&G(_,s),(!y||t&1)&&H!==(H=e[0].name+"")&&G(E,H);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const result = await pb.collection('${(ge=e[0])==null?void 0:ge.name}').listAuthMethods();
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final result = await pb.collection('${(we=e[0])==null?void 0:we.name}').listAuthMethods();
    `),P.$set(r),(!y||t&1)&&q!==(q=e[0].name+"")&&G(X,q),t&6&&(z=D(e[2]),g=ve(g,t,pe,1,e,z,he,O,je,Me,null,Pe)),t&6&&(N=D(e[2]),Ae(),b=ve(b,t,_e,1,e,N,fe,Q,Be,$e,null,Ce),Fe())},i(e){if(!y){W(P.$$.fragment,e),W(T.$$.fragment,e);for(let t=0;t<N.length;t+=1)W(b[t]);y=!0}},o(e){V(P.$$.fragment,e),V(T.$$.fragment,e);for(let t=0;t<b.length;t+=1)V(b[t]);y=!1},d(e){e&&(u(l),u(i),u(h),u(I),u(J),u(j),u(Y),u(M),u(Z),u(B),u(x),u($),u(le),u(F),u(oe),u(S)),ne(P,e),ne(T);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<b.length;t+=1)b[t].d()}}}function De(n,l,o){let s,{collection:_=new Oe}=l,f=200,i=[];const h=m=>o(1,f=m.code);return n.$$set=m=>{"collection"in m&&o(0,_=m.collection)},o(3,s=Qe.getApiExampleUrl(Ne.baseUrl)),o(2,i=[{code:200,body:`
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
            `}]),[_,f,i,s,h]}class Le extends Se{constructor(l){super(),ye(this,l,De,ze,Te,{collection:0})}}export{Le as default};
