import{S as re,i as ae,s as be,M as pe,C as P,e as p,w as y,b as a,c as te,f as u,g as t,h as I,m as ne,x as ue,t as ie,a as ce,o as n,d as le,T as me,p as de}from"./index-c8d873a4.js";import{S as fe}from"./SdkTabs-2977cee7.js";function $e(s){var B,U,W,T,A,H,L,M,q,j,J,N;let i,m,c=s[0].name+"",b,d,h,f,_,$,k,l,S,v,w,R,C,g,E,r,D;return l=new fe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${s[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(B=s[0])==null?void 0:B.name} record
        pb.collection('${(U=s[0])==null?void 0:U.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(W=s[0])==null?void 0:W.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(T=s[0])==null?void 0:T.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(A=s[0])==null?void 0:A.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(H=s[0])==null?void 0:H.name}').unsubscribe(); // remove all subscriptions in the collection
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${s[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(L=s[0])==null?void 0:L.name} record
        pb.collection('${(M=s[0])==null?void 0:M.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(q=s[0])==null?void 0:q.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        });

        // Unsubscribe
        pb.collection('${(j=s[0])==null?void 0:j.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(J=s[0])==null?void 0:J.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(N=s[0])==null?void 0:N.name}').unsubscribe(); // remove all subscriptions in the collection
    `}}),r=new pe({props:{content:JSON.stringify({action:"create",record:P.dummyCollectionRecord(s[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')}}),{c(){i=p("h3"),m=y("Realtime ("),b=y(c),d=y(")"),h=a(),f=p("div"),f.innerHTML=`<p>Subscribe to realtime changes via Server-Sent Events (SSE).</p> 
    <p>Events are sent for <strong>create</strong>, <strong>update</strong>
        and <strong>delete</strong> record operations (see &quot;Event data format&quot; section below).</p>`,_=a(),$=p("div"),$.innerHTML=`<div class="icon"><i class="ri-information-line"></i></div> 
    <div class="contet"><p><strong>You could subscribe to a single record or to an entire collection.</strong></p> 
        <p>When you subscribe to a <strong>single record</strong>, the collection&#39;s
            <strong>ViewRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p> 
        <p>When you subscribe to an <strong>entire collection</strong>, the collection&#39;s
            <strong>ListRule</strong> will be used to determine whether the subscriber has access to receive the
            event message.</p></div>`,k=a(),te(l.$$.fragment),S=a(),v=p("h6"),v.textContent="API details",w=a(),R=p("div"),R.innerHTML=`<strong class="label label-primary">SSE</strong> 
    <div class="content"><p>/api/realtime</p></div>`,C=a(),g=p("div"),g.textContent="Event data format",E=a(),te(r.$$.fragment),u(i,"class","m-b-sm"),u(f,"class","content txt-lg m-b-sm"),u($,"class","alert alert-info m-t-10 m-b-sm"),u(v,"class","m-b-xs"),u(R,"class","alert"),u(g,"class","section-title")},m(e,o){t(e,i,o),I(i,m),I(i,b),I(i,d),t(e,h,o),t(e,f,o),t(e,_,o),t(e,$,o),t(e,k,o),ne(l,e,o),t(e,S,o),t(e,v,o),t(e,w,o),t(e,R,o),t(e,C,o),t(e,g,o),t(e,E,o),ne(r,e,o),D=!0},p(e,[o]){var Y,z,F,G,K,Q,X,Z,x,ee,oe,se;(!D||o&1)&&c!==(c=e[0].name+"")&&ue(b,c);const O={};o&3&&(O.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(Y=e[0])==null?void 0:Y.name} record
        pb.collection('${(z=e[0])==null?void 0:z.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(F=e[0])==null?void 0:F.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        });

        // Unsubscribe
        pb.collection('${(G=e[0])==null?void 0:G.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(K=e[0])==null?void 0:K.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(Q=e[0])==null?void 0:Q.name}').unsubscribe(); // remove all subscriptions in the collection
    `),o&3&&(O.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[1]}');

        ...

        // (Optionally) authenticate
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // Subscribe to changes in any ${(X=e[0])==null?void 0:X.name} record
        pb.collection('${(Z=e[0])==null?void 0:Z.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        });

        // Subscribe to changes only in the specified record
        pb.collection('${(x=e[0])==null?void 0:x.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        });

        // Unsubscribe
        pb.collection('${(ee=e[0])==null?void 0:ee.name}').unsubscribe('RECORD_ID'); // remove all 'RECORD_ID' subscriptions
        pb.collection('${(oe=e[0])==null?void 0:oe.name}').unsubscribe('*'); // remove all '*' topic subscriptions
        pb.collection('${(se=e[0])==null?void 0:se.name}').unsubscribe(); // remove all subscriptions in the collection
    `),l.$set(O);const V={};o&1&&(V.content=JSON.stringify({action:"create",record:P.dummyCollectionRecord(e[0])},null,2).replace('"action": "create"','"action": "create" // create, update or delete')),r.$set(V)},i(e){D||(ie(l.$$.fragment,e),ie(r.$$.fragment,e),D=!0)},o(e){ce(l.$$.fragment,e),ce(r.$$.fragment,e),D=!1},d(e){e&&n(i),e&&n(h),e&&n(f),e&&n(_),e&&n($),e&&n(k),le(l,e),e&&n(S),e&&n(v),e&&n(w),e&&n(R),e&&n(C),e&&n(g),e&&n(E),le(r,e)}}}function ve(s,i,m){let c,{collection:b=new me}=i;return s.$$set=d=>{"collection"in d&&m(0,b=d.collection)},m(1,c=P.getApiExampleUrl(de.baseUrl)),[b,c]}class De extends re{constructor(i){super(),ae(this,i,ve,$e,be,{collection:0})}}export{De as default};
