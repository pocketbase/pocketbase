import{S as qe,i as Me,s as De,e as i,w as v,b as h,c as Se,f as m,g as d,h as s,m as Be,x as I,N as Te,P as Oe,k as We,Q as ze,n as He,t as le,a as oe,o as u,d as Ue,T as Le,C as je,p as Ie,r as N,u as Ne,M as Re}from"./index-a65ca895.js";import{S as Ke}from"./SdkTabs-ad912c8f.js";function ye(n,l,o){const a=n.slice();return a[5]=l[o],a}function Ae(n,l,o){const a=n.slice();return a[5]=l[o],a}function Ce(n,l){let o,a=l[5].code+"",_,b,c,p;function f(){return l[4](l[5])}return{key:n,first:null,c(){o=i("button"),_=v(a),b=h(),m(o,"class","tab-item"),N(o,"active",l[1]===l[5].code),this.first=o},m($,P){d($,o,P),s(o,_),s(o,b),c||(p=Ne(o,"click",f),c=!0)},p($,P){l=$,P&4&&a!==(a=l[5].code+"")&&I(_,a),P&6&&N(o,"active",l[1]===l[5].code)},d($){$&&u(o),c=!1,p()}}}function Ee(n,l){let o,a,_,b;return a=new Re({props:{content:l[5].body}}),{key:n,first:null,c(){o=i("div"),Se(a.$$.fragment),_=h(),m(o,"class","tab-item"),N(o,"active",l[1]===l[5].code),this.first=o},m(c,p){d(c,o,p),Be(a,o,null),s(o,_),b=!0},p(c,p){l=c;const f={};p&4&&(f.content=l[5].body),a.$set(f),(!b||p&6)&&N(o,"active",l[1]===l[5].code)},i(c){b||(le(a.$$.fragment,c),b=!0)},o(c){oe(a.$$.fragment,c),b=!1},d(c){c&&u(o),Ue(a)}}}function Qe(n){var he,_e,ke,ve;let l,o,a=n[0].name+"",_,b,c,p,f,$,P,O=n[0].name+"",R,se,ae,K,Q,y,F,E,G,w,W,ne,z,T,ie,J,H=n[0].name+"",V,ce,X,re,Y,de,L,Z,S,x,B,ee,U,te,A,q,g=[],ue=new Map,pe,M,k=[],fe=new Map,C;y=new Ke({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        await pb.collection('${(he=n[0])==null?void 0:he.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(_e=n[0])==null?void 0:_e.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        await pb.collection('${(ke=n[0])==null?void 0:ke.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(ve=n[0])==null?void 0:ve.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `}});let j=n[2];const me=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=Ae(n,j,e),r=me(t);ue.set(r,g[e]=Ce(r,t))}let D=n[2];const be=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=ye(n,D,e),r=be(t);fe.set(r,k[e]=Ee(r,t))}return{c(){l=i("h3"),o=v("Unlink OAuth2 account ("),_=v(a),b=v(")"),c=h(),p=i("div"),f=i("p"),$=v("Unlink a single external OAuth2 provider from "),P=i("strong"),R=v(O),se=v(" record."),ae=h(),K=i("p"),K.textContent="Only admins and the account owner can access this action.",Q=h(),Se(y.$$.fragment),F=h(),E=i("h6"),E.textContent="API details",G=h(),w=i("div"),W=i("strong"),W.textContent="DELETE",ne=h(),z=i("div"),T=i("p"),ie=v("/api/collections/"),J=i("strong"),V=v(H),ce=v("/records/"),X=i("strong"),X.textContent=":id",re=v("/external-auths/"),Y=i("strong"),Y.textContent=":provider",de=h(),L=i("p"),L.innerHTML="Requires <code>Authorization:TOKEN</code> header",Z=h(),S=i("div"),S.textContent="Path Parameters",x=h(),B=i("table"),B.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the auth record.</td></tr> 
        <tr><td>provider</td> 
            <td><span class="label">String</span></td> 
            <td>The name of the auth provider to unlink, eg. <code>google</code>, <code>twitter</code>,
                <code>github</code>, etc.</td></tr></tbody>`,ee=h(),U=i("div"),U.textContent="Responses",te=h(),A=i("div"),q=i("div");for(let e=0;e<g.length;e+=1)g[e].c();pe=h(),M=i("div");for(let e=0;e<k.length;e+=1)k[e].c();m(l,"class","m-b-sm"),m(p,"class","content txt-lg m-b-sm"),m(E,"class","m-b-xs"),m(W,"class","label label-primary"),m(z,"class","content"),m(L,"class","txt-hint txt-sm txt-right"),m(w,"class","alert alert-danger"),m(S,"class","section-title"),m(B,"class","table-compact table-border m-b-base"),m(U,"class","section-title"),m(q,"class","tabs-header compact left"),m(M,"class","tabs-content"),m(A,"class","tabs")},m(e,t){d(e,l,t),s(l,o),s(l,_),s(l,b),d(e,c,t),d(e,p,t),s(p,f),s(f,$),s(f,P),s(P,R),s(f,se),s(p,ae),s(p,K),d(e,Q,t),Be(y,e,t),d(e,F,t),d(e,E,t),d(e,G,t),d(e,w,t),s(w,W),s(w,ne),s(w,z),s(z,T),s(T,ie),s(T,J),s(J,V),s(T,ce),s(T,X),s(T,re),s(T,Y),s(w,de),s(w,L),d(e,Z,t),d(e,S,t),d(e,x,t),d(e,B,t),d(e,ee,t),d(e,U,t),d(e,te,t),d(e,A,t),s(A,q);for(let r=0;r<g.length;r+=1)g[r]&&g[r].m(q,null);s(A,pe),s(A,M);for(let r=0;r<k.length;r+=1)k[r]&&k[r].m(M,null);C=!0},p(e,[t]){var ge,we,$e,Pe;(!C||t&1)&&a!==(a=e[0].name+"")&&I(_,a),(!C||t&1)&&O!==(O=e[0].name+"")&&I(R,O);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ge=e[0])==null?void 0:ge.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(we=e[0])==null?void 0:we.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${($e=e[0])==null?void 0:$e.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(Pe=e[0])==null?void 0:Pe.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `),y.$set(r),(!C||t&1)&&H!==(H=e[0].name+"")&&I(V,H),t&6&&(j=e[2],g=Te(g,t,me,1,e,j,ue,q,Oe,Ce,null,Ae)),t&6&&(D=e[2],We(),k=Te(k,t,be,1,e,D,fe,M,ze,Ee,null,ye),He())},i(e){if(!C){le(y.$$.fragment,e);for(let t=0;t<D.length;t+=1)le(k[t]);C=!0}},o(e){oe(y.$$.fragment,e);for(let t=0;t<k.length;t+=1)oe(k[t]);C=!1},d(e){e&&u(l),e&&u(c),e&&u(p),e&&u(Q),Ue(y,e),e&&u(F),e&&u(E),e&&u(G),e&&u(w),e&&u(Z),e&&u(S),e&&u(x),e&&u(B),e&&u(ee),e&&u(U),e&&u(te),e&&u(A);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Fe(n,l,o){let a,{collection:_=new Le}=l,b=204,c=[];const p=f=>o(1,b=f.code);return n.$$set=f=>{"collection"in f&&o(0,_=f.collection)},o(3,a=je.getApiExampleUrl(Ie.baseUrl)),o(2,c=[{code:204,body:"null"},{code:401,body:`
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
            `},{code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}]),[_,b,c,a,p]}class Ve extends qe{constructor(l){super(),Me(this,l,Fe,Qe,De,{collection:0})}}export{Ve as default};
