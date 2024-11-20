import{S as xe,i as Ie,s as Je,U as Ke,V as Ne,W as J,f as s,y as k,h as p,c as K,j as b,l as d,n as o,m as Q,G as de,X as Le,Y as Qe,D as We,Z as Ge,E as Xe,t as V,a as U,u,d as W,I as Oe,p as Ye,k as G,o as Ze}from"./index-C6zau8vw.js";import{F as et}from"./FieldsQueryParam-Fbuu5Cus.js";function Ve(r,a,l){const n=r.slice();return n[5]=a[l],n}function Ue(r,a,l){const n=r.slice();return n[5]=a[l],n}function je(r,a){let l,n=a[5].code+"",m,_,i,h;function g(){return a[4](a[5])}return{key:r,first:null,c(){l=s("button"),m=k(n),_=p(),b(l,"class","tab-item"),G(l,"active",a[1]===a[5].code),this.first=l},m(v,w){d(v,l,w),o(l,m),o(l,_),i||(h=Ze(l,"click",g),i=!0)},p(v,w){a=v,w&4&&n!==(n=a[5].code+"")&&de(m,n),w&6&&G(l,"active",a[1]===a[5].code)},d(v){v&&u(l),i=!1,h()}}}function ze(r,a){let l,n,m,_;return n=new Ne({props:{content:a[5].body}}),{key:r,first:null,c(){l=s("div"),K(n.$$.fragment),m=p(),b(l,"class","tab-item"),G(l,"active",a[1]===a[5].code),this.first=l},m(i,h){d(i,l,h),Q(n,l,null),o(l,m),_=!0},p(i,h){a=i;const g={};h&4&&(g.content=a[5].body),n.$set(g),(!_||h&6)&&G(l,"active",a[1]===a[5].code)},i(i){_||(V(n.$$.fragment,i),_=!0)},o(i){U(n.$$.fragment,i),_=!1},d(i){i&&u(l),W(n)}}}function tt(r){var qe,Fe;let a,l,n=r[0].name+"",m,_,i,h,g,v,w,D,X,S,j,ue,z,M,pe,Y,N=r[0].name+"",Z,he,fe,x,ee,q,te,T,oe,be,F,C,ae,me,le,_e,f,ke,P,ge,ve,ye,se,$e,ne,Se,we,Te,re,Ce,Re,A,ie,E,ce,R,H,$=[],Pe=new Map,Ae,L,y=[],Be=new Map,B;v=new Ke({props:{js:`
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
    `}}),P=new Ne({props:{content:"?expand=relField1,relField2.subRelField"}}),A=new et({props:{prefix:"record."}});let I=J(r[2]);const De=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=Ue(r,I,e),c=De(t);Pe.set(c,$[e]=je(c,t))}let O=J(r[2]);const Me=e=>e[5].code;for(let e=0;e<O.length;e+=1){let t=Ve(r,O,e),c=Me(t);Be.set(c,y[e]=ze(c,t))}return{c(){a=s("h3"),l=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),h=s("div"),h.innerHTML=`<p>Returns a new auth response (token and record data) for an
        <strong>already authenticated record</strong>.</p> <p>This method is usually called by users on page/screen reload to ensure that the previously stored data
        in <code>pb.authStore</code> is still valid and up-to-date.</p>`,g=p(),K(v.$$.fragment),w=p(),D=s("h6"),D.textContent="API details",X=p(),S=s("div"),j=s("strong"),j.textContent="POST",ue=p(),z=s("div"),M=s("p"),pe=k("/api/collections/"),Y=s("strong"),Z=k(N),he=k("/auth-refresh"),fe=p(),x=s("p"),x.innerHTML="Requires <code>Authorization:TOKEN</code> header",ee=p(),q=s("div"),q.textContent="Query parameters",te=p(),T=s("table"),oe=s("thead"),oe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',be=p(),F=s("tbody"),C=s("tr"),ae=s("td"),ae.textContent="expand",me=p(),le=s("td"),le.innerHTML='<span class="label">String</span>',_e=p(),f=s("td"),ke=k(`Auto expand record relations. Ex.:
                `),K(P.$$.fragment),ge=k(`
                Supports up to 6-levels depth nested relations expansion. `),ve=s("br"),ye=k(`
                The expanded relations will be appended to the record under the
                `),se=s("code"),se.textContent="expand",$e=k(" property (eg. "),ne=s("code"),ne.textContent='"expand": {"relField1": {...}, ...}',Se=k(`).
                `),we=s("br"),Te=k(`
                Only the relations to which the request user has permissions to `),re=s("strong"),re.textContent="view",Ce=k(" will be expanded."),Re=p(),K(A.$$.fragment),ie=p(),E=s("div"),E.textContent="Responses",ce=p(),R=s("div"),H=s("div");for(let e=0;e<$.length;e+=1)$[e].c();Ae=p(),L=s("div");for(let e=0;e<y.length;e+=1)y[e].c();b(a,"class","m-b-sm"),b(h,"class","content txt-lg m-b-sm"),b(D,"class","m-b-xs"),b(j,"class","label label-primary"),b(z,"class","content"),b(x,"class","txt-hint txt-sm txt-right"),b(S,"class","alert alert-success"),b(q,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(E,"class","section-title"),b(H,"class","tabs-header compact combined left"),b(L,"class","tabs-content"),b(R,"class","tabs")},m(e,t){d(e,a,t),o(a,l),o(a,m),o(a,_),d(e,i,t),d(e,h,t),d(e,g,t),Q(v,e,t),d(e,w,t),d(e,D,t),d(e,X,t),d(e,S,t),o(S,j),o(S,ue),o(S,z),o(z,M),o(M,pe),o(M,Y),o(Y,Z),o(M,he),o(S,fe),o(S,x),d(e,ee,t),d(e,q,t),d(e,te,t),d(e,T,t),o(T,oe),o(T,be),o(T,F),o(F,C),o(C,ae),o(C,me),o(C,le),o(C,_e),o(C,f),o(f,ke),Q(P,f,null),o(f,ge),o(f,ve),o(f,ye),o(f,se),o(f,$e),o(f,ne),o(f,Se),o(f,we),o(f,Te),o(f,re),o(f,Ce),o(F,Re),Q(A,F,null),d(e,ie,t),d(e,E,t),d(e,ce,t),d(e,R,t),o(R,H);for(let c=0;c<$.length;c+=1)$[c]&&$[c].m(H,null);o(R,Ae),o(R,L);for(let c=0;c<y.length;c+=1)y[c]&&y[c].m(L,null);B=!0},p(e,[t]){var Ee,He;(!B||t&1)&&n!==(n=e[0].name+"")&&de(m,n);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(Ee=e[0])==null?void 0:Ee.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(He=e[0])==null?void 0:He.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);
    `),v.$set(c),(!B||t&1)&&N!==(N=e[0].name+"")&&de(Z,N),t&6&&(I=J(e[2]),$=Le($,t,De,1,e,I,Pe,H,Qe,je,null,Ue)),t&6&&(O=J(e[2]),We(),y=Le(y,t,Me,1,e,O,Be,L,Ge,ze,null,Ve),Xe())},i(e){if(!B){V(v.$$.fragment,e),V(P.$$.fragment,e),V(A.$$.fragment,e);for(let t=0;t<O.length;t+=1)V(y[t]);B=!0}},o(e){U(v.$$.fragment,e),U(P.$$.fragment,e),U(A.$$.fragment,e);for(let t=0;t<y.length;t+=1)U(y[t]);B=!1},d(e){e&&(u(a),u(i),u(h),u(g),u(w),u(D),u(X),u(S),u(ee),u(q),u(te),u(T),u(ie),u(E),u(ce),u(R)),W(v,e),W(P),W(A);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<y.length;t+=1)y[t].d()}}}function ot(r,a,l){let n,{collection:m}=a,_=200,i=[];const h=g=>l(1,_=g.code);return r.$$set=g=>{"collection"in g&&l(0,m=g.collection)},r.$$.update=()=>{r.$$.dirty&1&&l(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:Oe.dummyCollectionRecord(m)},null,2)},{code:401,body:`
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
            `}])},l(3,n=Oe.getApiExampleUrl(Ye.baseURL)),[m,_,i,n,h]}class st extends xe{constructor(a){super(),Ie(this,a,ot,tt,Je,{collection:0})}}export{st as default};
