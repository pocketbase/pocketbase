import{S as re,i as ae,s as be,N as ue,C as P,e as u,w as y,b as a,c as te,f as p,g as t,h as I,m as ne,x as pe,t as ie,a as le,o as n,d as ce,R as me,p as de}from"./index.89a3f554.js";import{S as fe}from"./SdkTabs.0a6ad1c9.js";function $e(o){var B,U,W,A,H,L,T,q,M,N,j,J;let i,m,l=o[0].name+"",b,d,h,f,_,$,k,c,S,v,w,R,C,g,E,r,D;return c=new fe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(B=o[0])==null?void 0:B.name} record
        pb.collection('${(U=o[0])==null?void 0:U.name}').subscribe('*', function (e) {
            console.log(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(W=o[0])==null?void 0:W.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(A=o[0])==null?void 0:A.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(H=o[0])==null?void 0:H.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(L=o[0])==null?void 0:L.name}').unsubscribe(); // remove all subscriptions in the collection
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(T=o[0])==null?void 0:T.name} record
        pb.collection('${(q=o[0])==null?void 0:q.name}').subscribe('*', (e) {
            console.log(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(M=o[0])==null?void 0:M.name}').subscribe('RECORD_ID', (e) {
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(N=o[0])==null?void 0:N.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(j=o[0])==null?void 0:j.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(J=o[0])==null?void 0:J.name}').unsubscribe(); // remove all subscriptions in the collection
    `}}),r=new ue({props:{content:JSON.stringify({action:"create",record:P.dummyCollectionRecord(o[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')}}),{c(){i=u("h3"),m=y("Realtime ("),b=y(l),d=y(")"),h=a(),f=u("div"),f.innerHTML=`<p>Subscribe to realtime changes via Server-Sent Events (SSE).</p> 
    <p>Events are sent for <strong>create</strong>, <strong>update</strong>
        and <strong>delete</strong> record operations (see &quot;Event data format&quot; section below).</p>`,_=a(),$=u("div"),$.innerHTML=`<div class="icon"><i class="ri-information-line"></i></div> 
    <div class="contet"><p><strong>You could subscribe to a single record or to an entire collection.</strong></p> 
        <p>When you subscribe to a <strong>single record</strong>, the collection&#39;s
            <strong>ViewRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p> 
        <p>When you subscribe to an <strong>entire collection</strong>, the collection&#39;s
            <strong>ListRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p></div>`,k=a(),te(c.$$.fragment),S=a(),v=u("h6"),v.textContent="API details",w=a(),R=u("div"),R.innerHTML=`<strong class="label label-primary">SSE</strong> 
    <div class="content"><p>/api/realtime</p></div>`,C=a(),g=u("div"),g.textContent="Event data format",E=a(),te(r.$$.fragment),p(i,"class","m-b-sm"),p(f,"class","content txt-lg m-b-sm"),p($,"class","alert alert-info m-t-10 m-b-sm"),p(v,"class","m-b-xs"),p(R,"class","alert"),p(g,"class","section-title")},m(e,s){t(e,i,s),I(i,m),I(i,b),I(i,d),t(e,h,s),t(e,f,s),t(e,_,s),t(e,$,s),t(e,k,s),ne(c,e,s),t(e,S,s),t(e,v,s),t(e,w,s),t(e,R,s),t(e,C,s),t(e,g,s),t(e,E,s),ne(r,e,s),D=!0},p(e,[s]){var Y,z,F,G,K,Q,X,Z,x,ee,se,oe;(!D||s&1)&&l!==(l=e[0].name+"")&&pe(b,l);const O={};s&3&&(O.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(Y=e[0])==null?void 0:Y.name} record
        pb.collection('${(z=e[0])==null?void 0:z.name}').subscribe('*', function (e) {
            console.log(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(F=e[0])==null?void 0:F.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(G=e[0])==null?void 0:G.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(K=e[0])==null?void 0:K.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(Q=e[0])==null?void 0:Q.name}').unsubscribe(); // remove all subscriptions in the collection
    `),s&3&&(O.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(X=e[0])==null?void 0:X.name} record
        pb.collection('${(Z=e[0])==null?void 0:Z.name}').subscribe('*', (e) {
            console.log(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(x=e[0])==null?void 0:x.name}').subscribe('RECORD_ID', (e) {
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(ee=e[0])==null?void 0:ee.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(se=e[0])==null?void 0:se.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(oe=e[0])==null?void 0:oe.name}').unsubscribe(); // remove all subscriptions in the collection
    `),c.$set(O);const V={};s&1&&(V.content=JSON.stringify({action:"create",record:P.dummyCollectionRecord(e[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')),r.$set(V)},i(e){D||(ie(c.$$.fragment,e),ie(r.$$.fragment,e),D=!0)},o(e){le(c.$$.fragment,e),le(r.$$.fragment,e),D=!1},d(e){e&&n(i),e&&n(h),e&&n(f),e&&n(_),e&&n($),e&&n(k),ce(c,e),e&&n(S),e&&n(v),e&&n(w),e&&n(R),e&&n(C),e&&n(g),e&&n(E),ce(r,e)}}}function ve(o,i,m){let l,{collection:b=new me}=i;return o.$$set=d=>{"collection"in d&&m(0,b=d.collection)},m(1,l=P.getApiExampleUrl(de.baseUrl)),[b,l]}class De extends re{constructor(i){super(),ae(this,i,ve,$e,be,{collection:0})}}export{De as default};
