import{S as Ue,i as xe,s as Ke,V as Qe,W as Ne,X as K,h as s,z as k,j as p,c as Q,k as b,n as d,o,m as W,H as de,Y as Le,Z as We,E as Ge,_ as Ie,G as Xe,t as V,a as z,v as u,d as G,J as Oe,p as Ye,l as I,q as Ze}from"./index-SKn09NMF.js";import{F as et}from"./FieldsQueryParam-NXRpAlxi.js";function Ve(r,a,l){const n=r.slice();return n[5]=a[l],n}function ze(r,a,l){const n=r.slice();return n[5]=a[l],n}function je(r,a){let l,n=a[5].code+"",m,_,i,h;function g(){return a[4](a[5])}return{key:r,first:null,c(){l=s("button"),m=k(n),_=p(),b(l,"class","tab-item"),I(l,"active",a[1]===a[5].code),this.first=l},m(v,w){d(v,l,w),o(l,m),o(l,_),i||(h=Ze(l,"click",g),i=!0)},p(v,w){a=v,w&4&&n!==(n=a[5].code+"")&&de(m,n),w&6&&I(l,"active",a[1]===a[5].code)},d(v){v&&u(l),i=!1,h()}}}function Je(r,a){let l,n,m,_;return n=new Ne({props:{content:a[5].body}}),{key:r,first:null,c(){l=s("div"),Q(n.$$.fragment),m=p(),b(l,"class","tab-item"),I(l,"active",a[1]===a[5].code),this.first=l},m(i,h){d(i,l,h),W(n,l,null),o(l,m),_=!0},p(i,h){a=i;const g={};h&4&&(g.content=a[5].body),n.$set(g),(!_||h&6)&&I(l,"active",a[1]===a[5].code)},i(i){_||(V(n.$$.fragment,i),_=!0)},o(i){z(n.$$.fragment,i),_=!1},d(i){i&&u(l),G(n)}}}function tt(r){var De,Fe;let a,l,n=r[0].name+"",m,_,i,h,g,v,w,B,X,S,j,ue,J,M,pe,Y,N=r[0].name+"",Z,he,fe,U,ee,D,te,T,oe,be,F,C,ae,me,le,_e,f,ke,P,ge,ve,$e,se,ye,ne,Se,we,Te,re,Ce,Re,q,ie,H,ce,R,E,y=[],Pe=new Map,qe,L,$=[],Ae=new Map,A;v=new Qe({props:{js:`
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
    `}}),P=new Ne({props:{content:"?expand=relField1,relField2.subRelField"}}),q=new et({props:{prefix:"record."}});let x=K(r[2]);const Be=e=>e[5].code;for(let e=0;e<x.length;e+=1){let t=ze(r,x,e),c=Be(t);Pe.set(c,y[e]=je(c,t))}let O=K(r[2]);const Me=e=>e[5].code;for(let e=0;e<O.length;e+=1){let t=Ve(r,O,e),c=Me(t);Ae.set(c,$[e]=Je(c,t))}return{c(){a=s("h3"),l=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),h=s("div"),h.innerHTML=`<p>Returns a new auth response (token and record data) for an
        <strong>already authenticated record</strong>.</p> <p>This method is usually called by users on page/screen reload to ensure that the previously stored data
        in <code>pb.authStore</code> is still valid and up-to-date.</p>`,g=p(),Q(v.$$.fragment),w=p(),B=s("h6"),B.textContent="API details",X=p(),S=s("div"),j=s("strong"),j.textContent="POST",ue=p(),J=s("div"),M=s("p"),pe=k("/api/collections/"),Y=s("strong"),Z=k(N),he=k("/auth-refresh"),fe=p(),U=s("p"),U.innerHTML="Requires <code>Authorization:TOKEN</code> header",ee=p(),D=s("div"),D.textContent="Query parameters",te=p(),T=s("table"),oe=s("thead"),oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',be=p(),F=s("tbody"),C=s("tr"),ae=s("td"),ae.textContent="expand",me=p(),le=s("td"),le.innerHTML='<span class="label">String</span>',_e=p(),f=s("td"),ke=k(`Auto expand record relations. Ex.:
                `),Q(P.$$.fragment),ge=k(`
                Supports up to 6-levels depth nested relations expansion. `),ve=s("br"),$e=k(`
                The expanded relations will be appended to the record under the
                `),se=s("code"),se.textContent="expand",ye=k(" property (eg. "),ne=s("code"),ne.textContent='"expand": {"relField1": {...}, ...}',Se=k(`).
                `),we=s("br"),Te=k(`
                Only the relations to which the request user has permissions to `),re=s("strong"),re.textContent="view",Ce=k(" will be expanded."),Re=p(),Q(q.$$.fragment),ie=p(),H=s("div"),H.textContent="Responses",ce=p(),R=s("div"),E=s("div");for(let e=0;e<y.length;e+=1)y[e].c();qe=p(),L=s("div");for(let e=0;e<$.length;e+=1)$[e].c();b(a,"class","m-b-sm"),b(h,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(j,"class","label label-primary"),b(J,"class","content"),b(U,"class","txt-hint txt-sm txt-right"),b(S,"class","alert alert-success"),b(D,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(H,"class","section-title"),b(E,"class","tabs-header compact combined left"),b(L,"class","tabs-content"),b(R,"class","tabs")},m(e,t){d(e,a,t),o(a,l),o(a,m),o(a,_),d(e,i,t),d(e,h,t),d(e,g,t),W(v,e,t),d(e,w,t),d(e,B,t),d(e,X,t),d(e,S,t),o(S,j),o(S,ue),o(S,J),o(J,M),o(M,pe),o(M,Y),o(Y,Z),o(M,he),o(S,fe),o(S,U),d(e,ee,t),d(e,D,t),d(e,te,t),d(e,T,t),o(T,oe),o(T,be),o(T,F),o(F,C),o(C,ae),o(C,me),o(C,le),o(C,_e),o(C,f),o(f,ke),W(P,f,null),o(f,ge),o(f,ve),o(f,$e),o(f,se),o(f,ye),o(f,ne),o(f,Se),o(f,we),o(f,Te),o(f,re),o(f,Ce),o(F,Re),W(q,F,null),d(e,ie,t),d(e,H,t),d(e,ce,t),d(e,R,t),o(R,E);for(let c=0;c<y.length;c+=1)y[c]&&y[c].m(E,null);o(R,qe),o(R,L);for(let c=0;c<$.length;c+=1)$[c]&&$[c].m(L,null);A=!0},p(e,[t]){var He,Ee;(!A||t&1)&&n!==(n=e[0].name+"")&&de(m,n);const c={};t&9&&(c.js=`
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

        final authData = await pb.collection('${(Ee=e[0])==null?void 0:Ee.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);
    `),v.$set(c),(!A||t&1)&&N!==(N=e[0].name+"")&&de(Z,N),t&6&&(x=K(e[2]),y=Le(y,t,Be,1,e,x,Pe,E,We,je,null,ze)),t&6&&(O=K(e[2]),Ge(),$=Le($,t,Me,1,e,O,Ae,L,Ie,Je,null,Ve),Xe())},i(e){if(!A){V(v.$$.fragment,e),V(P.$$.fragment,e),V(q.$$.fragment,e);for(let t=0;t<O.length;t+=1)V($[t]);A=!0}},o(e){z(v.$$.fragment,e),z(P.$$.fragment,e),z(q.$$.fragment,e);for(let t=0;t<$.length;t+=1)z($[t]);A=!1},d(e){e&&(u(a),u(i),u(h),u(g),u(w),u(B),u(X),u(S),u(ee),u(D),u(te),u(T),u(ie),u(H),u(ce),u(R)),G(v,e),G(P),G(q);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<$.length;t+=1)$[t].d()}}}function ot(r,a,l){let n,{collection:m}=a,_=200,i=[];const h=g=>l(1,_=g.code);return r.$$set=g=>{"collection"in g&&l(0,m=g.collection)},r.$$.update=()=>{r.$$.dirty&1&&l(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:Oe.dummyCollectionRecord(m)},null,2)},{code:401,body:`
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
            `}])},l(3,n=Oe.getApiExampleUrl(Ye.baseURL)),[m,_,i,n,h]}class st extends Ue{constructor(a){super(),xe(this,a,ot,tt,Ke,{collection:0})}}export{st as default};
