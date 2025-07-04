import{S as je,i as xe,s as Ie,a9 as Ke,aa as Ue,R as I,h as d,d as K,t as V,a as z,I as de,a0 as Oe,a1 as Qe,C as We,ab as Ge,D as Xe,l as u,n as a,m as Q,u as s,A as k,v as p,c as W,w as b,J as Ee,p as Ye,k as G,o as Ze}from"./index-CGZL2Fjn.js";import{F as et}from"./FieldsQueryParam-DQYOizuM.js";function Ve(r,o,l){const n=r.slice();return n[5]=o[l],n}function ze(r,o,l){const n=r.slice();return n[5]=o[l],n}function Je(r,o){let l,n=o[5].code+"",m,_,i,h;function g(){return o[4](o[5])}return{key:r,first:null,c(){l=s("button"),m=k(n),_=p(),b(l,"class","tab-item"),G(l,"active",o[1]===o[5].code),this.first=l},m(v,w){u(v,l,w),a(l,m),a(l,_),i||(h=Ze(l,"click",g),i=!0)},p(v,w){o=v,w&4&&n!==(n=o[5].code+"")&&de(m,n),w&6&&G(l,"active",o[1]===o[5].code)},d(v){v&&d(l),i=!1,h()}}}function Ne(r,o){let l,n,m,_;return n=new Ue({props:{content:o[5].body}}),{key:r,first:null,c(){l=s("div"),W(n.$$.fragment),m=p(),b(l,"class","tab-item"),G(l,"active",o[1]===o[5].code),this.first=l},m(i,h){u(i,l,h),Q(n,l,null),a(l,m),_=!0},p(i,h){o=i;const g={};h&4&&(g.content=o[5].body),n.$set(g),(!_||h&6)&&G(l,"active",o[1]===o[5].code)},i(i){_||(z(n.$$.fragment,i),_=!0)},o(i){V(n.$$.fragment,i),_=!1},d(i){i&&d(l),K(n)}}}function tt(r){var qe,Fe;let o,l,n=r[0].name+"",m,_,i,h,g,v,w,D,X,S,J,ue,N,M,pe,Y,U=r[0].name+"",Z,he,fe,j,ee,q,te,T,ae,be,F,C,oe,me,le,_e,f,ke,P,ge,ve,$e,se,ye,ne,Se,we,Te,re,Ce,Re,A,ie,H,ce,R,L,y=[],Pe=new Map,Ae,O,$=[],Be=new Map,B;v=new Ke({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${r[3]}');

        ...

        const authData = await pb.collection('${(qe=r[0])==null?void 0:qe.name}').authRefresh();

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
    `}}),P=new Ue({props:{content:"?expand=relField1,relField2.subRelField"}}),A=new et({props:{prefix:"record."}});let x=I(r[2]);const De=e=>e[5].code;for(let e=0;e<x.length;e+=1){let t=ze(r,x,e),c=De(t);Pe.set(c,y[e]=Je(c,t))}let E=I(r[2]);const Me=e=>e[5].code;for(let e=0;e<E.length;e+=1){let t=Ve(r,E,e),c=Me(t);Be.set(c,$[e]=Ne(c,t))}return{c(){o=s("h3"),l=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),h=s("div"),h.innerHTML=`<p>Returns a new auth response (token and record data) for an
        <strong>already authenticated record</strong>.</p> <p>This method is usually called by users on page/screen reload to ensure that the previously stored data
        in <code>pb.authStore</code> is still valid and up-to-date.</p>`,g=p(),W(v.$$.fragment),w=p(),D=s("h6"),D.textContent="API details",X=p(),S=s("div"),J=s("strong"),J.textContent="POST",ue=p(),N=s("div"),M=s("p"),pe=k("/api/collections/"),Y=s("strong"),Z=k(U),he=k("/auth-refresh"),fe=p(),j=s("p"),j.innerHTML="Requires <code>Authorization:TOKEN</code> header",ee=p(),q=s("div"),q.textContent="Query parameters",te=p(),T=s("table"),ae=s("thead"),ae.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',be=p(),F=s("tbody"),C=s("tr"),oe=s("td"),oe.textContent="expand",me=p(),le=s("td"),le.innerHTML='<span class="label">String</span>',_e=p(),f=s("td"),ke=k(`Auto expand record relations. Ex.:
                `),W(P.$$.fragment),ge=k(`
                Supports up to 6-levels depth nested relations expansion. `),ve=s("br"),$e=k(`
                The expanded relations will be appended to the record under the
                `),se=s("code"),se.textContent="expand",ye=k(" property (eg. "),ne=s("code"),ne.textContent='"expand": {"relField1": {...}, ...}',Se=k(`).
                `),we=s("br"),Te=k(`
                Only the relations to which the request user has permissions to `),re=s("strong"),re.textContent="view",Ce=k(" will be expanded."),Re=p(),W(A.$$.fragment),ie=p(),H=s("div"),H.textContent="Responses",ce=p(),R=s("div"),L=s("div");for(let e=0;e<y.length;e+=1)y[e].c();Ae=p(),O=s("div");for(let e=0;e<$.length;e+=1)$[e].c();b(o,"class","m-b-sm"),b(h,"class","content txt-lg m-b-sm"),b(D,"class","m-b-xs"),b(J,"class","label label-primary"),b(N,"class","content"),b(j,"class","txt-hint txt-sm txt-right"),b(S,"class","alert alert-success"),b(q,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(H,"class","section-title"),b(L,"class","tabs-header compact combined left"),b(O,"class","tabs-content"),b(R,"class","tabs")},m(e,t){u(e,o,t),a(o,l),a(o,m),a(o,_),u(e,i,t),u(e,h,t),u(e,g,t),Q(v,e,t),u(e,w,t),u(e,D,t),u(e,X,t),u(e,S,t),a(S,J),a(S,ue),a(S,N),a(N,M),a(M,pe),a(M,Y),a(Y,Z),a(M,he),a(S,fe),a(S,j),u(e,ee,t),u(e,q,t),u(e,te,t),u(e,T,t),a(T,ae),a(T,be),a(T,F),a(F,C),a(C,oe),a(C,me),a(C,le),a(C,_e),a(C,f),a(f,ke),Q(P,f,null),a(f,ge),a(f,ve),a(f,$e),a(f,se),a(f,ye),a(f,ne),a(f,Se),a(f,we),a(f,Te),a(f,re),a(f,Ce),a(F,Re),Q(A,F,null),u(e,ie,t),u(e,H,t),u(e,ce,t),u(e,R,t),a(R,L);for(let c=0;c<y.length;c+=1)y[c]&&y[c].m(L,null);a(R,Ae),a(R,O);for(let c=0;c<$.length;c+=1)$[c]&&$[c].m(O,null);B=!0},p(e,[t]){var He,Le;(!B||t&1)&&n!==(n=e[0].name+"")&&de(m,n);const c={};t&9&&(c.js=`
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
    `),v.$set(c),(!B||t&1)&&U!==(U=e[0].name+"")&&de(Z,U),t&6&&(x=I(e[2]),y=Oe(y,t,De,1,e,x,Pe,L,Qe,Je,null,ze)),t&6&&(E=I(e[2]),We(),$=Oe($,t,Me,1,e,E,Be,O,Ge,Ne,null,Ve),Xe())},i(e){if(!B){z(v.$$.fragment,e),z(P.$$.fragment,e),z(A.$$.fragment,e);for(let t=0;t<E.length;t+=1)z($[t]);B=!0}},o(e){V(v.$$.fragment,e),V(P.$$.fragment,e),V(A.$$.fragment,e);for(let t=0;t<$.length;t+=1)V($[t]);B=!1},d(e){e&&(d(o),d(i),d(h),d(g),d(w),d(D),d(X),d(S),d(ee),d(q),d(te),d(T),d(ie),d(H),d(ce),d(R)),K(v,e),K(P),K(A);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<$.length;t+=1)$[t].d()}}}function at(r,o,l){let n,{collection:m}=o,_=200,i=[];const h=g=>l(1,_=g.code);return r.$$set=g=>{"collection"in g&&l(0,m=g.collection)},r.$$.update=()=>{r.$$.dirty&1&&l(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:Ee.dummyCollectionRecord(m)},null,2)},{code:401,body:`
                {
                  "status": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `},{code:403,body:`
                {
                  "status": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `},{code:404,body:`
                {
                  "status": 404,
                  "message": "Missing auth record context.",
                  "data": {}
                }
            `}])},l(3,n=Ee.getApiExampleUrl(Ye.baseURL)),[m,_,i,n,h]}class st extends je{constructor(o){super(),xe(this,o,at,tt,Ie,{collection:0})}}export{st as default};
