import{S as St,i as At,s as Lt,V as Mt,W as Ht,X as Q,h as d,d as Re,t as Y,a as x,I as jt,Z as Pt,_ as Nt,C as Ut,$ as Jt,D as zt,l as u,n as t,m as Te,E as Wt,G as Gt,u as o,A as _,v as i,c as Pe,w as b,J as Ft,p as Kt,k as ee,o as Vt}from"./index-CkwOC79g.js";function Bt(a,s,n){const c=a.slice();return c[6]=s[n],c}function Et(a,s,n){const c=a.slice();return c[6]=s[n],c}function Ot(a,s){let n,c,y;function f(){return s[5](s[6])}return{key:a,first:null,c(){n=o("button"),n.textContent=`${s[6].code} `,b(n,"class","tab-item"),ee(n,"active",s[1]===s[6].code),this.first=n},m(r,h){u(r,n,h),c||(y=Vt(n,"click",f),c=!0)},p(r,h){s=r,h&10&&ee(n,"active",s[1]===s[6].code)},d(r){r&&d(n),c=!1,y()}}}function It(a,s){let n,c,y,f;return c=new Ht({props:{content:s[6].body}}),{key:a,first:null,c(){n=o("div"),Pe(c.$$.fragment),y=i(),b(n,"class","tab-item"),ee(n,"active",s[1]===s[6].code),this.first=n},m(r,h){u(r,n,h),Te(c,n,null),t(n,y),f=!0},p(r,h){s=r,(!f||h&10)&&ee(n,"active",s[1]===s[6].code)},i(r){f||(x(c.$$.fragment,r),f=!0)},o(r){Y(c.$$.fragment,r),f=!1},d(r){r&&d(n),Re(c)}}}function Xt(a){var pt,mt,bt,ht,ft,_t,yt,kt;let s,n,c=a[0].name+"",y,f,r,h,F,g,U,Fe,P,B,Be,E,Ee,Oe,te,le,w,oe,O,ae,I,se,H,ne,J,ie,q,ce,Ie,re,S,z,He,k,W,Se,de,Ae,D,G,Le,ue,Me,K,je,pe,Ne,C,Ue,me,Je,ze,We,V,Ge,X,Ke,be,Ve,he,Xe,fe,Ze,p,_e,Qe,ye,Ye,ke,xe,$e,et,ge,tt,ve,lt,ot,at,De,st,R,Ce,A,we,T,L,v=[],nt=new Map,it,M,$=[],ct=new Map,j,qe,rt;w=new Mt({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[2]}');

        ...

        const batch = pb.createBatch();

        batch.collection('${(pt=a[0])==null?void 0:pt.name}').create({ ... });
        batch.collection('${(mt=a[0])==null?void 0:mt.name}').update('RECORD_ID', { ... });
        batch.collection('${(bt=a[0])==null?void 0:bt.name}').delete('RECORD_ID');
        batch.collection('${(ht=a[0])==null?void 0:ht.name}').upsert({ ... });

        const result = await batch.send();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[2]}');

        ...

        final batch = pb.createBatch();

        batch.collection('${(ft=a[0])==null?void 0:ft.name}').create(body: { ... });
        batch.collection('${(_t=a[0])==null?void 0:_t.name}').update('RECORD_ID', body: { ... });
        batch.collection('${(yt=a[0])==null?void 0:yt.name}').delete('RECORD_ID');
        batch.collection('${(kt=a[0])==null?void 0:kt.name}').upsert(body: { ... });

        final result = await batch.send();
    `}}),R=new Ht({props:{language:"javascript",content:`
                            const formData = new FormData();

                            formData.append("@jsonPayload", JSON.stringify({
                                requests: [
                                    {
                                        method: "POST",
                                        url: "/api/collections/${a[0].name}/records?fields=id",
                                        body: { someField: "test1" }
                                    },
                                    {
                                        method: "PATCH",
                                        url: "/api/collections/${a[0].name}/records/RECORD_ID",
                                        body: { someField: "test2" }
                                    }
                                ]
                            }))

                            // file for the first request
                            formData.append("requests.0.someFileField", new File(...))

                            // file for the second request
                            formData.append("requests.1.someFileField", new File(...))
                        `}});let Z=Q(a[3]);const dt=e=>e[6].code;for(let e=0;e<Z.length;e+=1){let l=Et(a,Z,e),m=dt(l);nt.set(m,v[e]=Ot(m,l))}let N=Q(a[3]);const ut=e=>e[6].code;for(let e=0;e<N.length;e+=1){let l=Bt(a,N,e),m=ut(l);ct.set(m,$[e]=It(m,l))}return{c(){s=o("h3"),n=_("Batch create/update/upsert/delete ("),y=_(c),f=_(")"),r=i(),h=o("div"),h.innerHTML="<p>Batch and transactional create/update/upsert/delete of multiple records in a single request.</p>",F=i(),g=o("div"),U=o("div"),U.innerHTML='<i class="ri-error-warning-line"></i>',Fe=i(),P=o("div"),B=o("p"),Be=_(`The batch Web API need to be explicitly enabled and configured from the
            `),E=o("a"),E.textContent="Dashboard settings",Ee=_("."),Oe=i(),te=o("p"),te.innerHTML=`Because this endpoint process the requests in a single DB transaction it could degrade the
            performance of your application if not used with proper care and configuration
            <em>(prefer smaller max processing and body size limits, avoid large file uploads over slow S3
                networks and custom hooks that communicate with slow external APIs)</em>.`,le=i(),Pe(w.$$.fragment),oe=i(),O=o("h6"),O.textContent="API details",ae=i(),I=o("div"),I.innerHTML='<strong class="label label-primary">POST</strong> <div class="content">/api/batch</div>',se=i(),H=o("div"),H.textContent="Body Parameters",ne=i(),J=o("p"),J.innerHTML=`Body parameters could be sent as <em>application/json</em> or <em>multipart/form-data</em>.
    <br/>
    File upload is supported only via <em>multipart/form-data</em> (see below for more details).`,ie=i(),q=o("table"),ce=o("thead"),ce.innerHTML='<tr><th>Param</th> <th width="80%">Description</th></tr>',Ie=i(),re=o("tbody"),S=o("tr"),z=o("td"),z.innerHTML='<div class="flex txt-nowrap"><span class="label label-success">Required</span> <span>requests</span></div>',He=i(),k=o("td"),W=o("span"),W.textContent="Array<Request>",Se=_(` - List of the requests to process.

                `),de=o("p"),de.textContent="The supported batch request actions are:",Ae=i(),D=o("ul"),G=o("li"),Le=_("record create - "),ue=o("code"),ue.textContent="POST /api/collections/{collection}/records",Me=i(),K=o("li"),je=_(`record update -
                        `),pe=o("code"),pe.textContent="PATCH /api/collections/{collection}/records/{id}",Ne=i(),C=o("li"),Ue=_("record upsert - "),me=o("code"),me.textContent="PUT /api/collections/{collection}/records",Je=i(),ze=o("br"),We=i(),V=o("small"),V.innerHTML='(the body must have <code class="txt-sm">id</code> field)',Ge=i(),X=o("li"),Ke=_(`record delete -
                        `),be=o("code"),be.textContent="DELETE /api/collections/{collection}/records/{id}",Ve=i(),he=o("p"),he.textContent="Each batch Request element have the following properties:",Xe=i(),fe=o("ul"),fe.innerHTML=`<li><code>url path</code> <em>(could include query parameters)</em></li> <li><code>method</code> <em>(GET, POST, PUT, PATCH, DELETE)</em></li> <li><code>headers</code> <br/> <em>(custom per-request <code>Authorization</code> header is not supported at the moment,
                            aka. all batch requests have the same auth state)</em></li> <li><code>body</code></li>`,Ze=i(),p=o("p"),_e=o("strong"),_e.textContent="NB!",Qe=_(` When the batch request is send as
                    `),ye=o("code"),ye.textContent="multipart/form-data",Ye=_(`, the regular batch action fields are expected to be
                    submitted as serailized json under the `),ke=o("code"),ke.textContent="@jsonPayload",xe=_(` field and file keys need
                    to follow the pattern `),$e=o("code"),$e.textContent="requests.N.fileField",et=_(` or
                    `),ge=o("code"),ge.textContent="requests[N].fileField",tt=i(),ve=o("em"),ve.textContent=`(this is usually handled transparently by the SDKs when their specific object notation
                        is used)
                    `,lt=_(`.
                    `),ot=o("br"),at=_(`
                    If you don't use the SDKs or prefer manually to construct the `),De=o("code"),De.textContent="FormData",st=_(`
                    body, then it could look something like:
                    `),Pe(R.$$.fragment),Ce=i(),A=o("div"),A.textContent="Responses",we=i(),T=o("div"),L=o("div");for(let e=0;e<v.length;e+=1)v[e].c();it=i(),M=o("div");for(let e=0;e<$.length;e+=1)$[e].c();b(s,"class","m-b-sm"),b(h,"class","content txt-lg m-b-sm"),b(U,"class","icon"),b(E,"href","/settings"),b(P,"class","content"),b(g,"class","alert alert-warning"),b(O,"class","m-b-xs"),b(I,"class","api-route alert alert-success"),b(H,"class","section-title"),b(z,"valign","top"),b(W,"class","label"),b(V,"class","txt-hint"),b(q,"class","table-compact table-border m-t-xs m-b-base"),b(A,"class","section-title"),b(L,"class","tabs-header compact combined left"),b(M,"class","tabs-content"),b(T,"class","tabs")},m(e,l){u(e,s,l),t(s,n),t(s,y),t(s,f),u(e,r,l),u(e,h,l),u(e,F,l),u(e,g,l),t(g,U),t(g,Fe),t(g,P),t(P,B),t(B,Be),t(B,E),t(B,Ee),t(P,Oe),t(P,te),u(e,le,l),Te(w,e,l),u(e,oe,l),u(e,O,l),u(e,ae,l),u(e,I,l),u(e,se,l),u(e,H,l),u(e,ne,l),u(e,J,l),u(e,ie,l),u(e,q,l),t(q,ce),t(q,Ie),t(q,re),t(re,S),t(S,z),t(S,He),t(S,k),t(k,W),t(k,Se),t(k,de),t(k,Ae),t(k,D),t(D,G),t(G,Le),t(G,ue),t(D,Me),t(D,K),t(K,je),t(K,pe),t(D,Ne),t(D,C),t(C,Ue),t(C,me),t(C,Je),t(C,ze),t(C,We),t(C,V),t(D,Ge),t(D,X),t(X,Ke),t(X,be),t(k,Ve),t(k,he),t(k,Xe),t(k,fe),t(k,Ze),t(k,p),t(p,_e),t(p,Qe),t(p,ye),t(p,Ye),t(p,ke),t(p,xe),t(p,$e),t(p,et),t(p,ge),t(p,tt),t(p,ve),t(p,lt),t(p,ot),t(p,at),t(p,De),t(p,st),Te(R,p,null),u(e,Ce,l),u(e,A,l),u(e,we,l),u(e,T,l),t(T,L);for(let m=0;m<v.length;m+=1)v[m]&&v[m].m(L,null);t(T,it),t(T,M);for(let m=0;m<$.length;m+=1)$[m]&&$[m].m(M,null);j=!0,qe||(rt=Wt(Gt.call(null,E)),qe=!0)},p(e,[l]){var gt,vt,Dt,Ct,wt,qt,Rt,Tt;(!j||l&1)&&c!==(c=e[0].name+"")&&jt(y,c);const m={};l&5&&(m.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[2]}');

        ...

        const batch = pb.createBatch();

        batch.collection('${(gt=e[0])==null?void 0:gt.name}').create({ ... });
        batch.collection('${(vt=e[0])==null?void 0:vt.name}').update('RECORD_ID', { ... });
        batch.collection('${(Dt=e[0])==null?void 0:Dt.name}').delete('RECORD_ID');
        batch.collection('${(Ct=e[0])==null?void 0:Ct.name}').upsert({ ... });

        const result = await batch.send();
    `),l&5&&(m.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[2]}');

        ...

        final batch = pb.createBatch();

        batch.collection('${(wt=e[0])==null?void 0:wt.name}').create(body: { ... });
        batch.collection('${(qt=e[0])==null?void 0:qt.name}').update('RECORD_ID', body: { ... });
        batch.collection('${(Rt=e[0])==null?void 0:Rt.name}').delete('RECORD_ID');
        batch.collection('${(Tt=e[0])==null?void 0:Tt.name}').upsert(body: { ... });

        final result = await batch.send();
    `),w.$set(m);const $t={};l&1&&($t.content=`
                            const formData = new FormData();

                            formData.append("@jsonPayload", JSON.stringify({
                                requests: [
                                    {
                                        method: "POST",
                                        url: "/api/collections/${e[0].name}/records?fields=id",
                                        body: { someField: "test1" }
                                    },
                                    {
                                        method: "PATCH",
                                        url: "/api/collections/${e[0].name}/records/RECORD_ID",
                                        body: { someField: "test2" }
                                    }
                                ]
                            }))

                            // file for the first request
                            formData.append("requests.0.someFileField", new File(...))

                            // file for the second request
                            formData.append("requests.1.someFileField", new File(...))
                        `),R.$set($t),l&10&&(Z=Q(e[3]),v=Pt(v,l,dt,1,e,Z,nt,L,Nt,Ot,null,Et)),l&10&&(N=Q(e[3]),Ut(),$=Pt($,l,ut,1,e,N,ct,M,Jt,It,null,Bt),zt())},i(e){if(!j){x(w.$$.fragment,e),x(R.$$.fragment,e);for(let l=0;l<N.length;l+=1)x($[l]);j=!0}},o(e){Y(w.$$.fragment,e),Y(R.$$.fragment,e);for(let l=0;l<$.length;l+=1)Y($[l]);j=!1},d(e){e&&(d(s),d(r),d(h),d(F),d(g),d(le),d(oe),d(O),d(ae),d(I),d(se),d(H),d(ne),d(J),d(ie),d(q),d(Ce),d(A),d(we),d(T)),Re(w,e),Re(R);for(let l=0;l<v.length;l+=1)v[l].d();for(let l=0;l<$.length;l+=1)$[l].d();qe=!1,rt()}}}function Zt(a,s,n){let c,y,{collection:f}=s,r=200,h=[];const F=g=>n(1,r=g.code);return a.$$set=g=>{"collection"in g&&n(0,f=g.collection)},a.$$.update=()=>{a.$$.dirty&1&&n(4,y=Ft.dummyCollectionRecord(f)),a.$$.dirty&17&&f!=null&&f.id&&(h.push({code:200,body:JSON.stringify([{status:200,body:y},{status:200,body:Object.assign({},y,{id:y.id+"2"})}],null,2)}),h.push({code:400,body:`
                {
                  "status": 400,
                  "message": "Batch transaction failed.",
                  "data": {
                    "requests": {
                      "1": {
                        "code": "batch_request_failed",
                        "message": "Batch request failed.",
                        "response": {
                          "status": 400,
                          "message": "Failed to create record.",
                          "data": {
                            "id": {
                              "code": "validation_min_text_constraint",
                              "message": "Must be at least 3 character(s).",
                              "params": { "min": 3 }
                            }
                          }
                        }
                      }
                    }
                  }
                }
            `}),h.push({code:403,body:`
                {
                  "status": 403,
                  "message": "Batch requests are not allowed.",
                  "data": {}
                }
            `}))},n(2,c=Ft.getApiExampleUrl(Kt.baseURL)),[f,r,c,h,y,F]}class Yt extends St{constructor(s){super(),At(this,s,Zt,Xt,Lt,{collection:0})}}export{Yt as default};
