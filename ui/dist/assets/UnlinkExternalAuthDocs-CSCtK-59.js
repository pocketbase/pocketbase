import{S as Oe,i as De,s as Me,O as j,e as i,v as g,b as f,c as Be,f as h,g as d,h as a,m as qe,w as I,P as ye,Q as We,k as ze,R as He,n as Le,t as oe,a as ae,o as u,d as Ue,C as Re,A as je,q as N,r as Ie,N as Ne}from"./index-D0DO79Dq.js";import{S as Ke}from"./SdkTabs-DC6EUYpr.js";function Ce(n,l,o){const s=n.slice();return s[5]=l[o],s}function Te(n,l,o){const s=n.slice();return s[5]=l[o],s}function Ee(n,l){let o,s=l[5].code+"",_,b,c,p;function m(){return l[4](l[5])}return{key:n,first:null,c(){o=i("button"),_=g(s),b=f(),h(o,"class","tab-item"),N(o,"active",l[1]===l[5].code),this.first=o},m($,P){d($,o,P),a(o,_),a(o,b),c||(p=Ie(o,"click",m),c=!0)},p($,P){l=$,P&4&&s!==(s=l[5].code+"")&&I(_,s),P&6&&N(o,"active",l[1]===l[5].code)},d($){$&&u(o),c=!1,p()}}}function Se(n,l){let o,s,_,b;return s=new Ne({props:{content:l[5].body}}),{key:n,first:null,c(){o=i("div"),Be(s.$$.fragment),_=f(),h(o,"class","tab-item"),N(o,"active",l[1]===l[5].code),this.first=o},m(c,p){d(c,o,p),qe(s,o,null),a(o,_),b=!0},p(c,p){l=c;const m={};p&4&&(m.content=l[5].body),s.$set(m),(!b||p&6)&&N(o,"active",l[1]===l[5].code)},i(c){b||(oe(s.$$.fragment,c),b=!0)},o(c){ae(s.$$.fragment,c),b=!1},d(c){c&&u(o),Ue(s)}}}function Qe(n){var _e,ke,ge,ve;let l,o,s=n[0].name+"",_,b,c,p,m,$,P,M=n[0].name+"",K,se,ne,Q,F,y,G,E,J,w,W,ie,z,A,ce,V,H=n[0].name+"",X,re,Y,de,Z,ue,L,x,S,ee,B,te,q,le,C,U,v=[],pe=new Map,me,O,k=[],he=new Map,T;y=new Ke({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        await pb.collection('${(_e=n[0])==null?void 0:_e.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(ke=n[0])==null?void 0:ke.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        await pb.collection('${(ge=n[0])==null?void 0:ge.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(ve=n[0])==null?void 0:ve.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `}});let R=j(n[2]);const be=e=>e[5].code;for(let e=0;e<R.length;e+=1){let t=Te(n,R,e),r=be(t);pe.set(r,v[e]=Ee(r,t))}let D=j(n[2]);const fe=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=Ce(n,D,e),r=fe(t);he.set(r,k[e]=Se(r,t))}return{c(){l=i("h3"),o=g("Unlink OAuth2 account ("),_=g(s),b=g(")"),c=f(),p=i("div"),m=i("p"),$=g("Unlink a single external OAuth2 provider from "),P=i("strong"),K=g(M),se=g(" record."),ne=f(),Q=i("p"),Q.textContent="Only admins and the account owner can access this action.",F=f(),Be(y.$$.fragment),G=f(),E=i("h6"),E.textContent="API details",J=f(),w=i("div"),W=i("strong"),W.textContent="DELETE",ie=f(),z=i("div"),A=i("p"),ce=g("/api/collections/"),V=i("strong"),X=g(H),re=g("/records/"),Y=i("strong"),Y.textContent=":id",de=g("/external-auths/"),Z=i("strong"),Z.textContent=":provider",ue=f(),L=i("p"),L.innerHTML="Requires <code>Authorization:TOKEN</code> header",x=f(),S=i("div"),S.textContent="Path Parameters",ee=f(),B=i("table"),B.innerHTML=`<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the auth record.</td></tr> <tr><td>provider</td> <td><span class="label">String</span></td> <td>The name of the auth provider to unlink, eg. <code>google</code>, <code>twitter</code>,
                <code>github</code>, etc.</td></tr></tbody>`,te=f(),q=i("div"),q.textContent="Responses",le=f(),C=i("div"),U=i("div");for(let e=0;e<v.length;e+=1)v[e].c();me=f(),O=i("div");for(let e=0;e<k.length;e+=1)k[e].c();h(l,"class","m-b-sm"),h(p,"class","content txt-lg m-b-sm"),h(E,"class","m-b-xs"),h(W,"class","label label-primary"),h(z,"class","content"),h(L,"class","txt-hint txt-sm txt-right"),h(w,"class","alert alert-danger"),h(S,"class","section-title"),h(B,"class","table-compact table-border m-b-base"),h(q,"class","section-title"),h(U,"class","tabs-header compact combined left"),h(O,"class","tabs-content"),h(C,"class","tabs")},m(e,t){d(e,l,t),a(l,o),a(l,_),a(l,b),d(e,c,t),d(e,p,t),a(p,m),a(m,$),a(m,P),a(P,K),a(m,se),a(p,ne),a(p,Q),d(e,F,t),qe(y,e,t),d(e,G,t),d(e,E,t),d(e,J,t),d(e,w,t),a(w,W),a(w,ie),a(w,z),a(z,A),a(A,ce),a(A,V),a(V,X),a(A,re),a(A,Y),a(A,de),a(A,Z),a(w,ue),a(w,L),d(e,x,t),d(e,S,t),d(e,ee,t),d(e,B,t),d(e,te,t),d(e,q,t),d(e,le,t),d(e,C,t),a(C,U);for(let r=0;r<v.length;r+=1)v[r]&&v[r].m(U,null);a(C,me),a(C,O);for(let r=0;r<k.length;r+=1)k[r]&&k[r].m(O,null);T=!0},p(e,[t]){var we,$e,Pe,Ae;(!T||t&1)&&s!==(s=e[0].name+"")&&I(_,s),(!T||t&1)&&M!==(M=e[0].name+"")&&I(K,M);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(we=e[0])==null?void 0:we.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${($e=e[0])==null?void 0:$e.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(Pe=e[0])==null?void 0:Pe.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${(Ae=e[0])==null?void 0:Ae.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `),y.$set(r),(!T||t&1)&&H!==(H=e[0].name+"")&&I(X,H),t&6&&(R=j(e[2]),v=ye(v,t,be,1,e,R,pe,U,We,Ee,null,Te)),t&6&&(D=j(e[2]),ze(),k=ye(k,t,fe,1,e,D,he,O,He,Se,null,Ce),Le())},i(e){if(!T){oe(y.$$.fragment,e);for(let t=0;t<D.length;t+=1)oe(k[t]);T=!0}},o(e){ae(y.$$.fragment,e);for(let t=0;t<k.length;t+=1)ae(k[t]);T=!1},d(e){e&&(u(l),u(c),u(p),u(F),u(G),u(E),u(J),u(w),u(x),u(S),u(ee),u(B),u(te),u(q),u(le),u(C)),Ue(y,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Fe(n,l,o){let s,{collection:_}=l,b=204,c=[];const p=m=>o(1,b=m.code);return n.$$set=m=>{"collection"in m&&o(0,_=m.collection)},o(3,s=Re.getApiExampleUrl(je.baseUrl)),o(2,c=[{code:204,body:"null"},{code:401,body:`
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
            `}]),[_,b,c,s,p]}class Ve extends Oe{constructor(l){super(),De(this,l,Fe,Qe,Me,{collection:0})}}export{Ve as default};
