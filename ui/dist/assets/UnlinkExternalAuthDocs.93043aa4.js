import{S as qe,i as Oe,s as De,e as i,w as v,b as h,c as Se,f,g as r,h as s,m as Be,x as R,O as ye,P as Me,k as We,Q as ze,n as He,t as le,a as oe,o as d,d as Ue,R as Ie,C as Le,p as Re,r as j,u as je,N as Ne}from"./index.89a3f554.js";import{S as Ke}from"./SdkTabs.0a6ad1c9.js";function Ae(n,l,o){const a=n.slice();return a[5]=l[o],a}function Ce(n,l,o){const a=n.slice();return a[5]=l[o],a}function Te(n,l){let o,a=l[5].code+"",_,b,c,u;function m(){return l[4](l[5])}return{key:n,first:null,c(){o=i("button"),_=v(a),b=h(),f(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m($,P){r($,o,P),s(o,_),s(o,b),c||(u=je(o,"click",m),c=!0)},p($,P){l=$,P&4&&a!==(a=l[5].code+"")&&R(_,a),P&6&&j(o,"active",l[1]===l[5].code)},d($){$&&d(o),c=!1,u()}}}function Ee(n,l){let o,a,_,b;return a=new Ne({props:{content:l[5].body}}),{key:n,first:null,c(){o=i("div"),Se(a.$$.fragment),_=h(),f(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m(c,u){r(c,o,u),Be(a,o,null),s(o,_),b=!0},p(c,u){l=c;const m={};u&4&&(m.content=l[5].body),a.$set(m),(!b||u&6)&&j(o,"active",l[1]===l[5].code)},i(c){b||(le(a.$$.fragment,c),b=!0)},o(c){oe(a.$$.fragment,c),b=!1},d(c){c&&d(o),Ue(a)}}}function Qe(n){var he,_e,ke,ve;let l,o,a=n[0].name+"",_,b,c,u,m,$,P,M=n[0].name+"",N,se,ae,K,Q,A,F,E,G,g,W,ne,z,y,ie,J,H=n[0].name+"",V,ce,X,re,Y,de,I,Z,S,x,B,ee,U,te,C,q,w=[],ue=new Map,pe,O,k=[],me=new Map,T;A=new Ke({props:{js:`
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
    `}});let L=n[2];const fe=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=Ce(n,L,e),p=fe(t);ue.set(p,w[e]=Te(p,t))}let D=n[2];const be=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=Ae(n,D,e),p=be(t);me.set(p,k[e]=Ee(p,t))}return{c(){l=i("h3"),o=v("Unlink OAuth2 account ("),_=v(a),b=v(")"),c=h(),u=i("div"),m=i("p"),$=v("Unlink a single external OAuth2 provider from "),P=i("strong"),N=v(M),se=v(" record."),ae=h(),K=i("p"),K.textContent="Only admins and the account owner can access this action.",Q=h(),Se(A.$$.fragment),F=h(),E=i("h6"),E.textContent="API details",G=h(),g=i("div"),W=i("strong"),W.textContent="DELETE",ne=h(),z=i("div"),y=i("p"),ie=v("/api/collections/"),J=i("strong"),V=v(H),ce=v("/records/"),X=i("strong"),X.textContent=":id",re=v("/external-auths/"),Y=i("strong"),Y.textContent=":provider",de=h(),I=i("p"),I.innerHTML="Requires <code>Authorization:TOKEN</code> header",Z=h(),S=i("div"),S.textContent="Path Parameters",x=h(),B=i("table"),B.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the auth record.</td></tr> 
        <tr><td>provider</td> 
            <td><span class="label">String</span></td> 
            <td>The name of the auth provider to unlink, eg. <code>google</code>, <code>twitter</code>,
                <code>github</code>, etc.</td></tr></tbody>`,ee=h(),U=i("div"),U.textContent="Responses",te=h(),C=i("div"),q=i("div");for(let e=0;e<w.length;e+=1)w[e].c();pe=h(),O=i("div");for(let e=0;e<k.length;e+=1)k[e].c();f(l,"class","m-b-sm"),f(u,"class","content txt-lg m-b-sm"),f(E,"class","m-b-xs"),f(W,"class","label label-primary"),f(z,"class","content"),f(I,"class","txt-hint txt-sm txt-right"),f(g,"class","alert alert-danger"),f(S,"class","section-title"),f(B,"class","table-compact table-border m-b-base"),f(U,"class","section-title"),f(q,"class","tabs-header compact left"),f(O,"class","tabs-content"),f(C,"class","tabs")},m(e,t){r(e,l,t),s(l,o),s(l,_),s(l,b),r(e,c,t),r(e,u,t),s(u,m),s(m,$),s(m,P),s(P,N),s(m,se),s(u,ae),s(u,K),r(e,Q,t),Be(A,e,t),r(e,F,t),r(e,E,t),r(e,G,t),r(e,g,t),s(g,W),s(g,ne),s(g,z),s(z,y),s(y,ie),s(y,J),s(J,V),s(y,ce),s(y,X),s(y,re),s(y,Y),s(g,de),s(g,I),r(e,Z,t),r(e,S,t),r(e,x,t),r(e,B,t),r(e,ee,t),r(e,U,t),r(e,te,t),r(e,C,t),s(C,q);for(let p=0;p<w.length;p+=1)w[p].m(q,null);s(C,pe),s(C,O);for(let p=0;p<k.length;p+=1)k[p].m(O,null);T=!0},p(e,[t]){var ge,we,$e,Pe;(!T||t&1)&&a!==(a=e[0].name+"")&&R(_,a),(!T||t&1)&&M!==(M=e[0].name+"")&&R(N,M);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ge=e[0])==null?void 0:ge.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(we=e[0])==null?void 0:we.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${($e=e[0])==null?void 0:$e.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(Pe=e[0])==null?void 0:Pe.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `),A.$set(p),(!T||t&1)&&H!==(H=e[0].name+"")&&R(V,H),t&6&&(L=e[2],w=ye(w,t,fe,1,e,L,ue,q,Me,Te,null,Ce)),t&6&&(D=e[2],We(),k=ye(k,t,be,1,e,D,me,O,ze,Ee,null,Ae),He())},i(e){if(!T){le(A.$$.fragment,e);for(let t=0;t<D.length;t+=1)le(k[t]);T=!0}},o(e){oe(A.$$.fragment,e);for(let t=0;t<k.length;t+=1)oe(k[t]);T=!1},d(e){e&&d(l),e&&d(c),e&&d(u),e&&d(Q),Ue(A,e),e&&d(F),e&&d(E),e&&d(G),e&&d(g),e&&d(Z),e&&d(S),e&&d(x),e&&d(B),e&&d(ee),e&&d(U),e&&d(te),e&&d(C);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Fe(n,l,o){let a,{collection:_=new Ie}=l,b=204,c=[];const u=m=>o(1,b=m.code);return n.$$set=m=>{"collection"in m&&o(0,_=m.collection)},o(3,a=Le.getApiExampleUrl(Re.baseUrl)),o(2,c=[{code:204,body:"null"},{code:401,body:`
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
            `}]),[_,b,c,a,u]}class Ve extends qe{constructor(l){super(),Oe(this,l,Fe,Qe,De,{collection:0})}}export{Ve as default};
