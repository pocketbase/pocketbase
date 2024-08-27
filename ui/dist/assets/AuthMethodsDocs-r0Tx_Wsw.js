import{S as Se,i as ye,s as Ae,O as G,e as c,v as w,b as k,c as se,f as p,g as d,h as a,m as ae,w as U,P as ve,Q as Te,k as je,R as Be,n as Oe,t as W,a as V,o as u,d as ne,C as Fe,A as Qe,q as L,r as Ne,N as qe}from"./index-D0DO79Dq.js";import{S as He}from"./SdkTabs-DC6EUYpr.js";import{F as Ke}from"./FieldsQueryParam-BwleQAus.js";function Ce(n,l,o){const s=n.slice();return s[5]=l[o],s}function Pe(n,l,o){const s=n.slice();return s[5]=l[o],s}function $e(n,l){let o,s=l[5].code+"",_,f,i,h;function m(){return l[4](l[5])}return{key:n,first:null,c(){o=c("button"),_=w(s),f=k(),p(o,"class","tab-item"),L(o,"active",l[1]===l[5].code),this.first=o},m(v,C){d(v,o,C),a(o,_),a(o,f),i||(h=Ne(o,"click",m),i=!0)},p(v,C){l=v,C&4&&s!==(s=l[5].code+"")&&U(_,s),C&6&&L(o,"active",l[1]===l[5].code)},d(v){v&&u(o),i=!1,h()}}}function Me(n,l){let o,s,_,f;return s=new qe({props:{content:l[5].body}}),{key:n,first:null,c(){o=c("div"),se(s.$$.fragment),_=k(),p(o,"class","tab-item"),L(o,"active",l[1]===l[5].code),this.first=o},m(i,h){d(i,o,h),ae(s,o,null),a(o,_),f=!0},p(i,h){l=i;const m={};h&4&&(m.content=l[5].body),s.$set(m),(!f||h&6)&&L(o,"active",l[1]===l[5].code)},i(i){f||(W(s.$$.fragment,i),f=!0)},o(i){V(s.$$.fragment,i),f=!1},d(i){i&&u(o),ne(s)}}}function ze(n){var be,ke;let l,o,s=n[0].name+"",_,f,i,h,m,v,C,q=n[0].name+"",E,ie,I,P,J,T,Y,$,H,ce,K,j,re,R,z=n[0].name+"",X,de,Z,B,x,M,ee,ue,te,A,le,O,oe,S,F,g=[],he=new Map,me,Q,b=[],fe=new Map,y;P=new He({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        const result = await pb.collection('${(be=n[0])==null?void 0:be.name}').listAuthMethods();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        final result = await pb.collection('${(ke=n[0])==null?void 0:ke.name}').listAuthMethods();
    `}}),A=new Ke({});let D=G(n[2]);const pe=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=Pe(n,D,e),r=pe(t);he.set(r,g[e]=$e(r,t))}let N=G(n[2]);const _e=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=Ce(n,N,e),r=_e(t);fe.set(r,b[e]=Me(r,t))}return{c(){l=c("h3"),o=w("List auth methods ("),_=w(s),f=w(")"),i=k(),h=c("div"),m=c("p"),v=w("Returns a public list with all allowed "),C=c("strong"),E=w(q),ie=w(" authentication methods."),I=k(),se(P.$$.fragment),J=k(),T=c("h6"),T.textContent="API details",Y=k(),$=c("div"),H=c("strong"),H.textContent="GET",ce=k(),K=c("div"),j=c("p"),re=w("/api/collections/"),R=c("strong"),X=w(z),de=w("/auth-methods"),Z=k(),B=c("div"),B.textContent="Query parameters",x=k(),M=c("table"),ee=c("thead"),ee.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',ue=k(),te=c("tbody"),se(A.$$.fragment),le=k(),O=c("div"),O.textContent="Responses",oe=k(),S=c("div"),F=c("div");for(let e=0;e<g.length;e+=1)g[e].c();me=k(),Q=c("div");for(let e=0;e<b.length;e+=1)b[e].c();p(l,"class","m-b-sm"),p(h,"class","content txt-lg m-b-sm"),p(T,"class","m-b-xs"),p(H,"class","label label-primary"),p(K,"class","content"),p($,"class","alert alert-info"),p(B,"class","section-title"),p(M,"class","table-compact table-border m-b-base"),p(O,"class","section-title"),p(F,"class","tabs-header compact combined left"),p(Q,"class","tabs-content"),p(S,"class","tabs")},m(e,t){d(e,l,t),a(l,o),a(l,_),a(l,f),d(e,i,t),d(e,h,t),a(h,m),a(m,v),a(m,C),a(C,E),a(m,ie),d(e,I,t),ae(P,e,t),d(e,J,t),d(e,T,t),d(e,Y,t),d(e,$,t),a($,H),a($,ce),a($,K),a(K,j),a(j,re),a(j,R),a(R,X),a(j,de),d(e,Z,t),d(e,B,t),d(e,x,t),d(e,M,t),a(M,ee),a(M,ue),a(M,te),ae(A,te,null),d(e,le,t),d(e,O,t),d(e,oe,t),d(e,S,t),a(S,F);for(let r=0;r<g.length;r+=1)g[r]&&g[r].m(F,null);a(S,me),a(S,Q);for(let r=0;r<b.length;r+=1)b[r]&&b[r].m(Q,null);y=!0},p(e,[t]){var ge,we;(!y||t&1)&&s!==(s=e[0].name+"")&&U(_,s),(!y||t&1)&&q!==(q=e[0].name+"")&&U(E,q);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const result = await pb.collection('${(ge=e[0])==null?void 0:ge.name}').listAuthMethods();
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final result = await pb.collection('${(we=e[0])==null?void 0:we.name}').listAuthMethods();
    `),P.$set(r),(!y||t&1)&&z!==(z=e[0].name+"")&&U(X,z),t&6&&(D=G(e[2]),g=ve(g,t,pe,1,e,D,he,F,Te,$e,null,Pe)),t&6&&(N=G(e[2]),je(),b=ve(b,t,_e,1,e,N,fe,Q,Be,Me,null,Ce),Oe())},i(e){if(!y){W(P.$$.fragment,e),W(A.$$.fragment,e);for(let t=0;t<N.length;t+=1)W(b[t]);y=!0}},o(e){V(P.$$.fragment,e),V(A.$$.fragment,e);for(let t=0;t<b.length;t+=1)V(b[t]);y=!1},d(e){e&&(u(l),u(i),u(h),u(I),u(J),u(T),u(Y),u($),u(Z),u(B),u(x),u(M),u(le),u(O),u(oe),u(S)),ne(P,e),ne(A);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<b.length;t+=1)b[t].d()}}}function De(n,l,o){let s,{collection:_}=l,f=200,i=[];const h=m=>o(1,f=m.code);return n.$$set=m=>{"collection"in m&&o(0,_=m.collection)},o(3,s=Fe.getApiExampleUrl(Qe.baseUrl)),o(2,i=[{code:200,body:`
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
            `}]),[_,f,i,s,h]}class Ve extends Se{constructor(l){super(),ye(this,l,De,ze,Ae,{collection:0})}}export{Ve as default};
