import{S as je,i as xe,s as Je,N as Ue,O as J,e as s,v as k,b as p,c as K,f as b,g as d,h as o,m as I,w as de,P as Ee,Q as Ke,k as Ie,R as We,n as Ge,t as N,a as V,o as u,d as W,C as Le,A as Xe,q as G,r as Ye}from"./index-D0DO79Dq.js";import{S as Ze}from"./SdkTabs-DC6EUYpr.js";import{F as et}from"./FieldsQueryParam-BwleQAus.js";function Ne(r,l,a){const n=r.slice();return n[5]=l[a],n}function Ve(r,l,a){const n=r.slice();return n[5]=l[a],n}function ze(r,l){let a,n=l[5].code+"",m,_,i,h;function g(){return l[4](l[5])}return{key:r,first:null,c(){a=s("button"),m=k(n),_=p(),b(a,"class","tab-item"),G(a,"active",l[1]===l[5].code),this.first=a},m(v,w){d(v,a,w),o(a,m),o(a,_),i||(h=Ye(a,"click",g),i=!0)},p(v,w){l=v,w&4&&n!==(n=l[5].code+"")&&de(m,n),w&6&&G(a,"active",l[1]===l[5].code)},d(v){v&&u(a),i=!1,h()}}}function Qe(r,l){let a,n,m,_;return n=new Ue({props:{content:l[5].body}}),{key:r,first:null,c(){a=s("div"),K(n.$$.fragment),m=p(),b(a,"class","tab-item"),G(a,"active",l[1]===l[5].code),this.first=a},m(i,h){d(i,a,h),I(n,a,null),o(a,m),_=!0},p(i,h){l=i;const g={};h&4&&(g.content=l[5].body),n.$set(g),(!_||h&6)&&G(a,"active",l[1]===l[5].code)},i(i){_||(N(n.$$.fragment,i),_=!0)},o(i){V(n.$$.fragment,i),_=!1},d(i){i&&u(a),W(n)}}}function tt(r){var De,Fe;let l,a,n=r[0].name+"",m,_,i,h,g,v,w,B,X,S,z,ue,Q,M,pe,Y,U=r[0].name+"",Z,he,fe,j,ee,D,te,T,oe,be,F,C,le,me,ae,_e,f,ke,R,ge,ve,$e,se,ye,ne,Se,we,Te,re,Ce,Pe,A,ie,O,ce,P,H,y=[],Re=new Map,Ae,E,$=[],qe=new Map,q;v=new Ze({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${r[3]}');

        ...

        const authData = await pb.collection('${(De=r[0])==null?void 0:De.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${r[3]}');

        ...

        final authData = await pb.collection('${(Fe=r[0])==null?void 0:Fe.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);
    `}}),R=new Ue({props:{content:"?expand=relField1,relField2.subRelField"}}),A=new et({props:{prefix:"record."}});let x=J(r[2]);const Be=e=>e[5].code;for(let e=0;e<x.length;e+=1){let t=Ve(r,x,e),c=Be(t);Re.set(c,y[e]=ze(c,t))}let L=J(r[2]);const Me=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=Ne(r,L,e),c=Me(t);qe.set(c,$[e]=Qe(c,t))}return{c(){l=s("h3"),a=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),h=s("div"),h.innerHTML=`<p>Returns a new auth response (token and record data) for an
        <strong>already authenticated record</strong>.</p> <p><em>This method is usually called by users on page/screen reload to ensure that the previously stored
            data in <code>pb.authStore</code> is still valid and up-to-date.</em></p>`,g=p(),K(v.$$.fragment),w=p(),B=s("h6"),B.textContent="API details",X=p(),S=s("div"),z=s("strong"),z.textContent="POST",ue=p(),Q=s("div"),M=s("p"),pe=k("/api/collections/"),Y=s("strong"),Z=k(U),he=k("/auth-refresh"),fe=p(),j=s("p"),j.innerHTML="Requires record <code>Authorization:TOKEN</code> header",ee=p(),D=s("div"),D.textContent="Query parameters",te=p(),T=s("table"),oe=s("thead"),oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',be=p(),F=s("tbody"),C=s("tr"),le=s("td"),le.textContent="expand",me=p(),ae=s("td"),ae.innerHTML='<span class="label">String</span>',_e=p(),f=s("td"),ke=k(`Auto expand record relations. Ex.:
                `),K(R.$$.fragment),ge=k(`
                Supports up to 6-levels depth nested relations expansion. `),ve=s("br"),$e=k(`
                The expanded relations will be appended to the record under the
                `),se=s("code"),se.textContent="expand",ye=k(" property (eg. "),ne=s("code"),ne.textContent='"expand": {"relField1": {...}, ...}',Se=k(`).
                `),we=s("br"),Te=k(`
                Only the relations to which the request user has permissions to `),re=s("strong"),re.textContent="view",Ce=k(" will be expanded."),Pe=p(),K(A.$$.fragment),ie=p(),O=s("div"),O.textContent="Responses",ce=p(),P=s("div"),H=s("div");for(let e=0;e<y.length;e+=1)y[e].c();Ae=p(),E=s("div");for(let e=0;e<$.length;e+=1)$[e].c();b(l,"class","m-b-sm"),b(h,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(z,"class","label label-primary"),b(Q,"class","content"),b(j,"class","txt-hint txt-sm txt-right"),b(S,"class","alert alert-success"),b(D,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(O,"class","section-title"),b(H,"class","tabs-header compact combined left"),b(E,"class","tabs-content"),b(P,"class","tabs")},m(e,t){d(e,l,t),o(l,a),o(l,m),o(l,_),d(e,i,t),d(e,h,t),d(e,g,t),I(v,e,t),d(e,w,t),d(e,B,t),d(e,X,t),d(e,S,t),o(S,z),o(S,ue),o(S,Q),o(Q,M),o(M,pe),o(M,Y),o(Y,Z),o(M,he),o(S,fe),o(S,j),d(e,ee,t),d(e,D,t),d(e,te,t),d(e,T,t),o(T,oe),o(T,be),o(T,F),o(F,C),o(C,le),o(C,me),o(C,ae),o(C,_e),o(C,f),o(f,ke),I(R,f,null),o(f,ge),o(f,ve),o(f,$e),o(f,se),o(f,ye),o(f,ne),o(f,Se),o(f,we),o(f,Te),o(f,re),o(f,Ce),o(F,Pe),I(A,F,null),d(e,ie,t),d(e,O,t),d(e,ce,t),d(e,P,t),o(P,H);for(let c=0;c<y.length;c+=1)y[c]&&y[c].m(H,null);o(P,Ae),o(P,E);for(let c=0;c<$.length;c+=1)$[c]&&$[c].m(E,null);q=!0},p(e,[t]){var Oe,He;(!q||t&1)&&n!==(n=e[0].name+"")&&de(m,n);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(Oe=e[0])==null?void 0:Oe.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(He=e[0])==null?void 0:He.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);
    `),v.$set(c),(!q||t&1)&&U!==(U=e[0].name+"")&&de(Z,U),t&6&&(x=J(e[2]),y=Ee(y,t,Be,1,e,x,Re,H,Ke,ze,null,Ve)),t&6&&(L=J(e[2]),Ie(),$=Ee($,t,Me,1,e,L,qe,E,We,Qe,null,Ne),Ge())},i(e){if(!q){N(v.$$.fragment,e),N(R.$$.fragment,e),N(A.$$.fragment,e);for(let t=0;t<L.length;t+=1)N($[t]);q=!0}},o(e){V(v.$$.fragment,e),V(R.$$.fragment,e),V(A.$$.fragment,e);for(let t=0;t<$.length;t+=1)V($[t]);q=!1},d(e){e&&(u(l),u(i),u(h),u(g),u(w),u(B),u(X),u(S),u(ee),u(D),u(te),u(T),u(ie),u(O),u(ce),u(P)),W(v,e),W(R),W(A);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<$.length;t+=1)$[t].d()}}}function ot(r,l,a){let n,{collection:m}=l,_=200,i=[];const h=g=>a(1,_=g.code);return r.$$set=g=>{"collection"in g&&a(0,m=g.collection)},r.$$.update=()=>{r.$$.dirty&1&&a(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:Le.dummyCollectionRecord(m)},null,2)},{code:401,body:`
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
            `}])},a(3,n=Le.getApiExampleUrl(Xe.baseUrl)),[m,_,i,n,h]}class nt extends je{constructor(l){super(),xe(this,l,ot,tt,Je,{collection:0})}}export{nt as default};
