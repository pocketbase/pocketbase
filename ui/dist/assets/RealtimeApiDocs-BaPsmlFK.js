import{S as re,i as ae,s as be,Q as pe,R as ue,C as P,e as p,w as y,b as a,c as se,f as u,g as s,h as I,m as ne,x as me,t as ie,a as ce,o as n,d as le,p as de}from"./index-DlBS34em.js";function he(o){var B,U,W,L,A,H,T,q,M,j,J,N;let i,m,c=o[0].name+"",b,d,k,h,D,f,_,l,C,$,S,g,w,v,E,r,R;return l=new pe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(B=o[0])==null?void 0:B.name} record
        pb.collection('${(U=o[0])==null?void 0:U.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like: filter, expand, custom headers, etc. */ });

        // Subscribe to changes only in the specified record
        pb.collection('${(W=o[0])==null?void 0:W.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like: filter, expand, custom headers, etc. */ });

        // Unsubscribe
        pb.collection('${(L=o[0])==null?void 0:L.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(A=o[0])==null?void 0:A.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(H=o[0])==null?void 0:H.name}').unsubscribe(); // remove all subscriptions in the collection
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(T=o[0])==null?void 0:T.name} record
        pb.collection('${(q=o[0])==null?void 0:q.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like: filter, expand, custom headers, etc. */);

        // Subscribe to changes only in the specified record
        pb.collection('${(M=o[0])==null?void 0:M.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like: filter, expand, custom headers, etc. */);

        // Unsubscribe
        pb.collection('${(j=o[0])==null?void 0:j.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(J=o[0])==null?void 0:J.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(N=o[0])==null?void 0:N.name}').unsubscribe(); // remove all subscriptions in the collection
    `}}),r=new ue({props:{content:JSON.stringify({action:"create",record:P.dummyCollectionRecord(o[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')}}),{c(){i=p("h3"),m=y("Realtime ("),b=y(c),d=y(")"),k=a(),h=p("div"),h.innerHTML=`<p>Subscribe to realtime changes via Server-Sent Events (SSE).</p> <p>Events are sent for <strong>create</strong>, <strong>update</strong>
        and <strong>delete</strong> record operations (see &quot;Event data format&quot; section below).</p>`,D=a(),f=p("div"),f.innerHTML=`<div class="icon"><i class="ri-information-line"></i></div> <div class="contet"><p><strong>You could subscribe to a single record or to an entire collection.</strong></p> <p>When you subscribe to a <strong>single record</strong>, the collection&#39;s
            <strong>ViewRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p> <p>When you subscribe to an <strong>entire collection</strong>, the collection&#39;s
            <strong>ListRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p></div>`,_=a(),se(l.$$.fragment),C=a(),$=p("h6"),$.textContent="API details",S=a(),g=p("div"),g.innerHTML='<strong class="label label-primary">SSE</strong> <div class="content"><p>/api/realtime</p></div>',w=a(),v=p("div"),v.textContent="Event data format",E=a(),se(r.$$.fragment),u(i,"class","m-b-sm"),u(h,"class","content txt-lg m-b-sm"),u(f,"class","alert alert-info m-t-10 m-b-sm"),u($,"class","m-b-xs"),u(g,"class","alert"),u(v,"class","section-title")},m(e,t){s(e,i,t),I(i,m),I(i,b),I(i,d),s(e,k,t),s(e,h,t),s(e,D,t),s(e,f,t),s(e,_,t),ne(l,e,t),s(e,C,t),s(e,$,t),s(e,S,t),s(e,g,t),s(e,w,t),s(e,v,t),s(e,E,t),ne(r,e,t),R=!0},p(e,[t]){var V,Y,z,F,G,K,X,Z,x,ee,te,oe;(!R||t&1)&&c!==(c=e[0].name+"")&&me(b,c);const O={};t&3&&(O.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(V=e[0])==null?void 0:V.name} record
        pb.collection('${(Y=e[0])==null?void 0:Y.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like: filter, expand, custom headers, etc. */ });

        // Subscribe to changes only in the specified record
        pb.collection('${(z=e[0])==null?void 0:z.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* other options like: filter, expand, custom headers, etc. */ });

        // Unsubscribe
        pb.collection('${(F=e[0])==null?void 0:F.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(G=e[0])==null?void 0:G.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(K=e[0])==null?void 0:K.name}').unsubscribe(); // remove all subscriptions in the collection
    `),t&3&&(O.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(X=e[0])==null?void 0:X.name} record
        pb.collection('${(Z=e[0])==null?void 0:Z.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like: filter, expand, custom headers, etc. */);

        // Subscribe to changes only in the specified record
        pb.collection('${(x=e[0])==null?void 0:x.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        }, /* other options like: filter, expand, custom headers, etc. */);

        // Unsubscribe
        pb.collection('${(ee=e[0])==null?void 0:ee.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(te=e[0])==null?void 0:te.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(oe=e[0])==null?void 0:oe.name}').unsubscribe(); // remove all subscriptions in the collection
    `),l.$set(O);const Q={};t&1&&(Q.content=JSON.stringify({action:"create",record:P.dummyCollectionRecord(e[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')),r.$set(Q)},i(e){R||(ie(l.$$.fragment,e),ie(r.$$.fragment,e),R=!0)},o(e){ce(l.$$.fragment,e),ce(r.$$.fragment,e),R=!1},d(e){e&&(n(i),n(k),n(h),n(D),n(f),n(_),n(C),n($),n(S),n(g),n(w),n(v),n(E)),le(l,e),le(r,e)}}}function fe(o,i,m){let c,{collection:b}=i;return o.$$set=d=>{"collection"in d&&m(0,b=d.collection)},m(1,c=P.getApiExampleUrl(de.baseURL)),[b,c]}class ge extends re{constructor(i){super(),ae(this,i,fe,he,be,{collection:0})}}export{ge as default};
