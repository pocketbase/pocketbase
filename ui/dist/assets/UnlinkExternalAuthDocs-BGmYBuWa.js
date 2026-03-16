import{S as Oe,i as De,s as Me,O as j,b as d,d as Be,t as oe,a as ae,r as I,Q as Te,R as We,g as ze,T as He,e as Le,f as u,h as a,m as Ue,n as i,u as g,k as f,c as qe,o as b,C as Re,p as je,w as N,x as Ie,N as Ne}from"./index-lKVVd1Bs.js";import{S as Ke}from"./SdkTabs-CROZp_fs.js";function Ae(n,l,o){const s=n.slice();return s[5]=l[o],s}function Ce(n,l,o){const s=n.slice();return s[5]=l[o],s}function Ee(n,l){let o,s=l[5].code+"",_,h,c,p;function m(){return l[4](l[5])}return{key:n,first:null,c(){o=i("button"),_=g(s),h=f(),b(o,"class","tab-item"),N(o,"active",l[1]===l[5].code),this.first=o},m($,P){u($,o,P),a(o,_),a(o,h),c||(p=Ie(o,"click",m),c=!0)},p($,P){l=$,P&4&&s!==(s=l[5].code+"")&&I(_,s),P&6&&N(o,"active",l[1]===l[5].code)},d($){$&&d(o),c=!1,p()}}}function Se(n,l){let o,s,_,h;return s=new Ne({props:{content:l[5].body}}),{key:n,first:null,c(){o=i("div"),qe(s.$$.fragment),_=f(),b(o,"class","tab-item"),N(o,"active",l[1]===l[5].code),this.first=o},m(c,p){u(c,o,p),Ue(s,o,null),a(o,_),h=!0},p(c,p){l=c;const m={};p&4&&(m.content=l[5].body),s.$set(m),(!h||p&6)&&N(o,"active",l[1]===l[5].code)},i(c){h||(ae(s.$$.fragment,c),h=!0)},o(c){oe(s.$$.fragment,c),h=!1},d(c){c&&d(o),Be(s)}}}function Qe(n){var _e,ke,ge,ve;let l,o,s=n[0].name+"",_,h,c,p,m,$,P,M=n[0].name+"",K,se,ne,Q,F,T,G,E,J,w,W,ie,z,y,ce,V,H=n[0].name+"",X,re,Y,de,Z,ue,L,x,S,ee,B,te,U,le,A,q,v=[],pe=new Map,me,O,k=[],be=new Map,C;T=new Ke({props:{js:`
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
    `}});let R=j(n[2]);const he=e=>e[5].code;for(let e=0;e<R.length;e+=1){let t=Ce(n,R,e),r=he(t);pe.set(r,v[e]=Ee(r,t))}let D=j(n[2]);const fe=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=Ae(n,D,e),r=fe(t);be.set(r,k[e]=Se(r,t))}return{c(){l=i("h3"),o=g("Unlink OAuth2 account ("),_=g(s),h=g(")"),c=f(),p=i("div"),m=i("p"),$=g("Unlink a single external OAuth2 provider from "),P=i("strong"),K=g(M),se=g(" record."),ne=f(),Q=i("p"),Q.textContent="Only admins and the account owner can access this action.",F=f(),qe(T.$$.fragment),G=f(),E=i("h6"),E.textContent="API details",J=f(),w=i("div"),W=i("strong"),W.textContent="DELETE",ie=f(),z=i("div"),y=i("p"),ce=g("/api/collections/"),V=i("strong"),X=g(H),re=g("/records/"),Y=i("strong"),Y.textContent=":id",de=g("/external-auths/"),Z=i("strong"),Z.textContent=":provider",ue=f(),L=i("p"),L.innerHTML="Requires <code>Authorization:TOKEN</code> header",x=f(),S=i("div"),S.textContent="Path Parameters",ee=f(),B=i("table"),B.innerHTML=`<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the auth record.</td></tr> <tr><td>provider</td> <td><span class="label">String</span></td> <td>The name of the auth provider to unlink, eg. <code>google</code>, <code>twitter</code>,
                <code>github</code>, etc.</td></tr></tbody>`,te=f(),U=i("div"),U.textContent="Responses",le=f(),A=i("div"),q=i("div");for(let e=0;e<v.length;e+=1)v[e].c();me=f(),O=i("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(E,"class","m-b-xs"),b(W,"class","label label-primary"),b(z,"class","content"),b(L,"class","txt-hint txt-sm txt-right"),b(w,"class","alert alert-danger"),b(S,"class","section-title"),b(B,"class","table-compact table-border m-b-base"),b(U,"class","section-title"),b(q,"class","tabs-header compact combined left"),b(O,"class","tabs-content"),b(A,"class","tabs")},m(e,t){u(e,l,t),a(l,o),a(l,_),a(l,h),u(e,c,t),u(e,p,t),a(p,m),a(m,$),a(m,P),a(P,K),a(m,se),a(p,ne),a(p,Q),u(e,F,t),Ue(T,e,t),u(e,G,t),u(e,E,t),u(e,J,t),u(e,w,t),a(w,W),a(w,ie),a(w,z),a(z,y),a(y,ce),a(y,V),a(V,X),a(y,re),a(y,Y),a(y,de),a(y,Z),a(w,ue),a(w,L),u(e,x,t),u(e,S,t),u(e,ee,t),u(e,B,t),u(e,te,t),u(e,U,t),u(e,le,t),u(e,A,t),a(A,q);for(let r=0;r<v.length;r+=1)v[r]&&v[r].m(q,null);a(A,me),a(A,O);for(let r=0;r<k.length;r+=1)k[r]&&k[r].m(O,null);C=!0},p(e,[t]){var we,$e,Pe,ye;(!C||t&1)&&s!==(s=e[0].name+"")&&I(_,s),(!C||t&1)&&M!==(M=e[0].name+"")&&I(K,M);const r={};t&9&&(r.js=`
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

        await pb.collection('${(ye=e[0])==null?void 0:ye.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `),T.$set(r),(!C||t&1)&&H!==(H=e[0].name+"")&&I(X,H),t&6&&(R=j(e[2]),v=Te(v,t,he,1,e,R,pe,q,We,Ee,null,Ce)),t&6&&(D=j(e[2]),ze(),k=Te(k,t,fe,1,e,D,be,O,He,Se,null,Ae),Le())},i(e){if(!C){ae(T.$$.fragment,e);for(let t=0;t<D.length;t+=1)ae(k[t]);C=!0}},o(e){oe(T.$$.fragment,e);for(let t=0;t<k.length;t+=1)oe(k[t]);C=!1},d(e){e&&(d(l),d(c),d(p),d(F),d(G),d(E),d(J),d(w),d(x),d(S),d(ee),d(B),d(te),d(U),d(le),d(A)),Be(T,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Fe(n,l,o){let s,{collection:_}=l,h=204,c=[];const p=m=>o(1,h=m.code);return n.$$set=m=>{"collection"in m&&o(0,_=m.collection)},o(3,s=Re.getApiExampleUrl(je.baseUrl)),o(2,c=[{code:204,body:"null"},{code:401,body:`
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
            `}]),[_,h,c,s,p]}class Ve extends Oe{constructor(l){super(),De(this,l,Fe,Qe,Me,{collection:0})}}export{Ve as default};
