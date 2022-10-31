import{S as qe,i as Oe,s as De,e as i,w as v,b as h,c as Se,f,g as r,h as s,m as Be,x as R,P as ye,Q as Le,k as Me,R as Ve,n as ze,t as le,a as oe,o as d,d as Ue,L as He,C as Ie,p as Re,r as j,u as je,O as Ke}from"./index.7b2502cb.js";import{S as Ne}from"./SdkTabs.315f7f19.js";function Ae(n,l,o){const a=n.slice();return a[5]=l[o],a}function Ce(n,l,o){const a=n.slice();return a[5]=l[o],a}function Pe(n,l){let o,a=l[5].code+"",_,b,c,u;function p(){return l[4](l[5])}return{key:n,first:null,c(){o=i("button"),_=v(a),b=h(),f(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m($,E){r($,o,E),s(o,_),s(o,b),c||(u=je(o,"click",p),c=!0)},p($,E){l=$,E&4&&a!==(a=l[5].code+"")&&R(_,a),E&6&&j(o,"active",l[1]===l[5].code)},d($){$&&d(o),c=!1,u()}}}function Te(n,l){let o,a,_,b;return a=new Ke({props:{content:l[5].body}}),{key:n,first:null,c(){o=i("div"),Se(a.$$.fragment),_=h(),f(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m(c,u){r(c,o,u),Be(a,o,null),s(o,_),b=!0},p(c,u){l=c;const p={};u&4&&(p.content=l[5].body),a.$set(p),(!b||u&6)&&j(o,"active",l[1]===l[5].code)},i(c){b||(le(a.$$.fragment,c),b=!0)},o(c){oe(a.$$.fragment,c),b=!1},d(c){c&&d(o),Ue(a)}}}function Qe(n){var he,_e,ke,ve;let l,o,a=n[0].name+"",_,b,c,u,p,$,E,L=n[0].name+"",K,se,ae,N,Q,A,F,T,G,g,M,ne,V,y,ie,J,z=n[0].name+"",W,ce,X,re,Y,de,H,Z,S,x,B,ee,U,te,C,q,w=[],ue=new Map,me,O,k=[],pe=new Map,P;A=new Ne({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        await pb.collection('${(he=n[0])==null?void 0:he.name}').authViaEmail('test@example.com', '123456');

        await pb.collection('${(_e=n[0])==null?void 0:_e.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        await pb.collection('${(ke=n[0])==null?void 0:ke.name}').authViaEmail('test@example.com', '123456');

        await pb.collection('${(ve=n[0])==null?void 0:ve.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `}});let I=n[2];const fe=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=Ce(n,I,e),m=fe(t);ue.set(m,w[e]=Pe(m,t))}let D=n[2];const be=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=Ae(n,D,e),m=be(t);pe.set(m,k[e]=Te(m,t))}return{c(){l=i("h3"),o=v("Unlink OAuth2 account ("),_=v(a),b=v(")"),c=h(),u=i("div"),p=i("p"),$=v("Unlink a single external OAuth2 provider from "),E=i("strong"),K=v(L),se=v(" record."),ae=h(),N=i("p"),N.textContent="Only admins and the account owner can access this action.",Q=h(),Se(A.$$.fragment),F=h(),T=i("h6"),T.textContent="API details",G=h(),g=i("div"),M=i("strong"),M.textContent="DELETE",ne=h(),V=i("div"),y=i("p"),ie=v("/api/collections/"),J=i("strong"),W=v(z),ce=v("/records/"),X=i("strong"),X.textContent=":id",re=v("/external-auths/"),Y=i("strong"),Y.textContent=":provider",de=h(),H=i("p"),H.innerHTML="Requires <code>Authorization:TOKEN</code> header",Z=h(),S=i("div"),S.textContent="Path Parameters",x=h(),B=i("table"),B.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the auth record.</td></tr> 
        <tr><td>provider</td> 
            <td><span class="label">String</span></td> 
            <td>The name of the auth provider to unlink, eg. <code>google</code>, <code>twitter</code>,
                <code>github</code>, etc.</td></tr></tbody>`,ee=h(),U=i("div"),U.textContent="Responses",te=h(),C=i("div"),q=i("div");for(let e=0;e<w.length;e+=1)w[e].c();me=h(),O=i("div");for(let e=0;e<k.length;e+=1)k[e].c();f(l,"class","m-b-sm"),f(u,"class","content txt-lg m-b-sm"),f(T,"class","m-b-xs"),f(M,"class","label label-primary"),f(V,"class","content"),f(H,"class","txt-hint txt-sm txt-right"),f(g,"class","alert alert-danger"),f(S,"class","section-title"),f(B,"class","table-compact table-border m-b-base"),f(U,"class","section-title"),f(q,"class","tabs-header compact left"),f(O,"class","tabs-content"),f(C,"class","tabs")},m(e,t){r(e,l,t),s(l,o),s(l,_),s(l,b),r(e,c,t),r(e,u,t),s(u,p),s(p,$),s(p,E),s(E,K),s(p,se),s(u,ae),s(u,N),r(e,Q,t),Be(A,e,t),r(e,F,t),r(e,T,t),r(e,G,t),r(e,g,t),s(g,M),s(g,ne),s(g,V),s(V,y),s(y,ie),s(y,J),s(J,W),s(y,ce),s(y,X),s(y,re),s(y,Y),s(g,de),s(g,H),r(e,Z,t),r(e,S,t),r(e,x,t),r(e,B,t),r(e,ee,t),r(e,U,t),r(e,te,t),r(e,C,t),s(C,q);for(let m=0;m<w.length;m+=1)w[m].m(q,null);s(C,me),s(C,O);for(let m=0;m<k.length;m+=1)k[m].m(O,null);P=!0},p(e,[t]){var ge,we,$e,Ee;(!P||t&1)&&a!==(a=e[0].name+"")&&R(_,a),(!P||t&1)&&L!==(L=e[0].name+"")&&R(K,L);const m={};t&9&&(m.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ge=e[0])==null?void 0:ge.name}').authViaEmail('test@example.com', '123456');

        await pb.collection('${(we=e[0])==null?void 0:we.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `),t&9&&(m.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${($e=e[0])==null?void 0:$e.name}').authViaEmail('test@example.com', '123456');

        await pb.collection('${(Ee=e[0])==null?void 0:Ee.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `),A.$set(m),(!P||t&1)&&z!==(z=e[0].name+"")&&R(W,z),t&6&&(I=e[2],w=ye(w,t,fe,1,e,I,ue,q,Le,Pe,null,Ce)),t&6&&(D=e[2],Me(),k=ye(k,t,be,1,e,D,pe,O,Ve,Te,null,Ae),ze())},i(e){if(!P){le(A.$$.fragment,e);for(let t=0;t<D.length;t+=1)le(k[t]);P=!0}},o(e){oe(A.$$.fragment,e);for(let t=0;t<k.length;t+=1)oe(k[t]);P=!1},d(e){e&&d(l),e&&d(c),e&&d(u),e&&d(Q),Ue(A,e),e&&d(F),e&&d(T),e&&d(G),e&&d(g),e&&d(Z),e&&d(S),e&&d(x),e&&d(B),e&&d(ee),e&&d(U),e&&d(te),e&&d(C);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Fe(n,l,o){let a,{collection:_=new He}=l,b=204,c=[];const u=p=>o(1,b=p.code);return n.$$set=p=>{"collection"in p&&o(0,_=p.collection)},o(3,a=Ie.getApiExampleUrl(Re.baseUrl)),o(2,c=[{code:204,body:"null"},{code:401,body:`
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
            `}]),[_,b,c,a,u]}class We extends qe{constructor(l){super(),Oe(this,l,Fe,Qe,De,{collection:0})}}export{We as default};
