import{S as Qe,i as je,s as Je,Q as Ke,R as Ne,T as J,e as s,w as k,b as p,c as K,f as b,g as d,h as o,m as W,x as de,U as Oe,V as We,k as Ie,W as Ge,n as Xe,t as E,a as U,o as u,d as I,C as Ve,p as Ye,r as G,u as Ze}from"./index-DuKqYKLn.js";import{F as et}from"./FieldsQueryParam-JQgEaE1u.js";function Ee(r,a,l){const n=r.slice();return n[5]=a[l],n}function Ue(r,a,l){const n=r.slice();return n[5]=a[l],n}function xe(r,a){let l,n=a[5].code+"",m,_,i,h;function g(){return a[4](a[5])}return{key:r,first:null,c(){l=s("button"),m=k(n),_=p(),b(l,"class","tab-item"),G(l,"active",a[1]===a[5].code),this.first=l},m(v,w){d(v,l,w),o(l,m),o(l,_),i||(h=Ze(l,"click",g),i=!0)},p(v,w){a=v,w&4&&n!==(n=a[5].code+"")&&de(m,n),w&6&&G(l,"active",a[1]===a[5].code)},d(v){v&&u(l),i=!1,h()}}}function ze(r,a){let l,n,m,_;return n=new Ne({props:{content:a[5].body}}),{key:r,first:null,c(){l=s("div"),K(n.$$.fragment),m=p(),b(l,"class","tab-item"),G(l,"active",a[1]===a[5].code),this.first=l},m(i,h){d(i,l,h),W(n,l,null),o(l,m),_=!0},p(i,h){a=i;const g={};h&4&&(g.content=a[5].body),n.$set(g),(!_||h&6)&&G(l,"active",a[1]===a[5].code)},i(i){_||(E(n.$$.fragment,i),_=!0)},o(i){U(n.$$.fragment,i),_=!1},d(i){i&&u(l),I(n)}}}function tt(r){var De,Fe;let a,l,n=r[0].name+"",m,_,i,h,g,v,w,M,X,S,x,ue,z,q,pe,Y,N=r[0].name+"",Z,he,fe,Q,ee,D,te,T,oe,be,F,C,ae,me,le,_e,f,ke,P,ge,ve,$e,se,ye,ne,Se,we,Te,re,Ce,Re,A,ie,H,ce,R,L,y=[],Pe=new Map,Ae,O,$=[],Be=new Map,B;v=new Ke({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${r[3]}');

        ...

        const authData = await pb.collection('${(De=r[0])==null?void 0:De.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${r[3]}');

        ...

        final authData = await pb.collection('${(Fe=r[0])==null?void 0:Fe.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);
    `}}),P=new Ne({props:{content:"?expand=relField1,relField2.subRelField"}}),A=new et({props:{prefix:"record."}});let j=J(r[2]);const Me=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=Ue(r,j,e),c=Me(t);Pe.set(c,y[e]=xe(c,t))}let V=J(r[2]);const qe=e=>e[5].code;for(let e=0;e<V.length;e+=1){let t=Ee(r,V,e),c=qe(t);Be.set(c,$[e]=ze(c,t))}return{c(){a=s("h3"),l=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),h=s("div"),h.innerHTML=`<p>Returns a new auth response (token and record data) for an
        <strong>already authenticated record</strong>.</p> <p>This method is usually called by users on page/screen reload to ensure that the previously stored data
        in <code>pb.authStore</code> is still valid and up-to-date.</p>`,g=p(),K(v.$$.fragment),w=p(),M=s("h6"),M.textContent="API details",X=p(),S=s("div"),x=s("strong"),x.textContent="POST",ue=p(),z=s("div"),q=s("p"),pe=k("/api/collections/"),Y=s("strong"),Z=k(N),he=k("/auth-refresh"),fe=p(),Q=s("p"),Q.innerHTML="Requires <code>Authorization:TOKEN</code> header",ee=p(),D=s("div"),D.textContent="Query parameters",te=p(),T=s("table"),oe=s("thead"),oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',be=p(),F=s("tbody"),C=s("tr"),ae=s("td"),ae.textContent="expand",me=p(),le=s("td"),le.innerHTML='<span class="label">String</span>',_e=p(),f=s("td"),ke=k(`Auto expand record relations. Ex.:
                `),K(P.$$.fragment),ge=k(`
                Supports up to 6-levels depth nested relations expansion. `),ve=s("br"),$e=k(`
                The expanded relations will be appended to the record under the
                `),se=s("code"),se.textContent="expand",ye=k(" property (eg. "),ne=s("code"),ne.textContent='"expand": {"relField1": {...}, ...}',Se=k(`).
                `),we=s("br"),Te=k(`
                Only the relations to which the request user has permissions to `),re=s("strong"),re.textContent="view",Ce=k(" will be expanded."),Re=p(),K(A.$$.fragment),ie=p(),H=s("div"),H.textContent="Responses",ce=p(),R=s("div"),L=s("div");for(let e=0;e<y.length;e+=1)y[e].c();Ae=p(),O=s("div");for(let e=0;e<$.length;e+=1)$[e].c();b(a,"class","m-b-sm"),b(h,"class","content txt-lg m-b-sm"),b(M,"class","m-b-xs"),b(x,"class","label label-primary"),b(z,"class","content"),b(Q,"class","txt-hint txt-sm txt-right"),b(S,"class","alert alert-success"),b(D,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(H,"class","section-title"),b(L,"class","tabs-header compact combined left"),b(O,"class","tabs-content"),b(R,"class","tabs")},m(e,t){d(e,a,t),o(a,l),o(a,m),o(a,_),d(e,i,t),d(e,h,t),d(e,g,t),W(v,e,t),d(e,w,t),d(e,M,t),d(e,X,t),d(e,S,t),o(S,x),o(S,ue),o(S,z),o(z,q),o(q,pe),o(q,Y),o(Y,Z),o(q,he),o(S,fe),o(S,Q),d(e,ee,t),d(e,D,t),d(e,te,t),d(e,T,t),o(T,oe),o(T,be),o(T,F),o(F,C),o(C,ae),o(C,me),o(C,le),o(C,_e),o(C,f),o(f,ke),W(P,f,null),o(f,ge),o(f,ve),o(f,$e),o(f,se),o(f,ye),o(f,ne),o(f,Se),o(f,we),o(f,Te),o(f,re),o(f,Ce),o(F,Re),W(A,F,null),d(e,ie,t),d(e,H,t),d(e,ce,t),d(e,R,t),o(R,L);for(let c=0;c<y.length;c+=1)y[c]&&y[c].m(L,null);o(R,Ae),o(R,O);for(let c=0;c<$.length;c+=1)$[c]&&$[c].m(O,null);B=!0},p(e,[t]){var He,Le;(!B||t&1)&&n!==(n=e[0].name+"")&&de(m,n);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(He=e[0])==null?void 0:He.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(Le=e[0])==null?void 0:Le.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);
    `),v.$set(c),(!B||t&1)&&N!==(N=e[0].name+"")&&de(Z,N),t&6&&(j=J(e[2]),y=Oe(y,t,Me,1,e,j,Pe,L,We,xe,null,Ue)),t&6&&(V=J(e[2]),Ie(),$=Oe($,t,qe,1,e,V,Be,O,Ge,ze,null,Ee),Xe())},i(e){if(!B){E(v.$$.fragment,e),E(P.$$.fragment,e),E(A.$$.fragment,e);for(let t=0;t<V.length;t+=1)E($[t]);B=!0}},o(e){U(v.$$.fragment,e),U(P.$$.fragment,e),U(A.$$.fragment,e);for(let t=0;t<$.length;t+=1)U($[t]);B=!1},d(e){e&&(u(a),u(i),u(h),u(g),u(w),u(M),u(X),u(S),u(ee),u(D),u(te),u(T),u(ie),u(H),u(ce),u(R)),I(v,e),I(P),I(A);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<$.length;t+=1)$[t].d()}}}function ot(r,a,l){let n,{collection:m}=a,_=200,i=[];const h=g=>l(1,_=g.code);return r.$$set=g=>{"collection"in g&&l(0,m=g.collection)},r.$$.update=()=>{r.$$.dirty&1&&l(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:Ve.dummyCollectionRecord(m)},null,2)},{code:401,body:`
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
                  "message": "Missing auth record context.",
                  "data": {}
                }
            `}])},l(3,n=Ve.getApiExampleUrl(Ye.baseURL)),[m,_,i,n,h]}class st extends Qe{constructor(a){super(),je(this,a,ot,tt,Je,{collection:0})}}export{st as default};
