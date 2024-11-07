import{S as Ht,i as Lt,s as jt,U as At,V as St,W as Z,f as o,y as _,h as i,c as Re,j as b,l as d,n as t,m as Te,B as Mt,C as Nt,G as Ut,X as Pt,Y as zt,D as Jt,Z as Wt,E as Gt,t as Q,a as x,u,d as Pe,I as Ft,p as Kt,k as ee,o as Vt}from"./index-D8ZreWy5.js";function Bt(a,s,n){const c=a.slice();return c[6]=s[n],c}function Et(a,s,n){const c=a.slice();return c[6]=s[n],c}function Ot(a,s){let n,c,y;function h(){return s[5](s[6])}return{key:a,first:null,c(){n=o("button"),n.textContent=`${s[6].code} `,b(n,"class","tab-item"),ee(n,"active",s[1]===s[6].code),this.first=n},m(r,f){d(r,n,f),c||(y=Vt(n,"click",h),c=!0)},p(r,f){s=r,f&10&&ee(n,"active",s[1]===s[6].code)},d(r){r&&u(n),c=!1,y()}}}function It(a,s){let n,c,y,h;return c=new St({props:{content:s[6].body}}),{key:a,first:null,c(){n=o("div"),Re(c.$$.fragment),y=i(),b(n,"class","tab-item"),ee(n,"active",s[1]===s[6].code),this.first=n},m(r,f){d(r,n,f),Te(c,n,null),t(n,y),h=!0},p(r,f){s=r,(!h||f&10)&&ee(n,"active",s[1]===s[6].code)},i(r){h||(Q(c.$$.fragment,r),h=!0)},o(r){x(c.$$.fragment,r),h=!1},d(r){r&&u(n),Pe(c)}}}function Xt(a){var pt,mt,bt,ft,ht,_t,yt,kt;let s,n,c=a[0].name+"",y,h,r,f,F,g,U,Fe,P,B,Be,E,Ee,Oe,te,le,w,oe,O,ae,I,se,S,ne,z,ie,q,ce,Ie,re,H,J,Se,k,W,He,de,Le,D,G,je,ue,Ae,K,Me,pe,Ne,v,Ue,me,ze,Je,We,V,Ge,X,Ke,be,Ve,fe,Xe,he,Ye,p,_e,Ze,ye,Qe,ke,xe,$e,et,ge,tt,Ce,lt,ot,at,De,st,R,ve,L,we,T,j,C=[],nt=new Map,it,A,$=[],ct=new Map,M,qe,rt;w=new At({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[2]}');

        ...

        const batch = pb.createBatch();

        batch.collection('${(pt=a[0])==null?void 0:pt.name}').create({ ... });
        batch.collection('${(mt=a[0])==null?void 0:mt.name}').update("RECORD_ID", { ... });
        batch.collection('${(bt=a[0])==null?void 0:bt.name}').delete("RECORD_ID");
        batch.collection('${(ft=a[0])==null?void 0:ft.name}').upsert({ ... });

        const result = await batch.send();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[2]}');

        ...

        final batch = pb.createBatch();

        batch.collection('${(ht=a[0])==null?void 0:ht.name}').create(body: { ... });
        batch.collection('${(_t=a[0])==null?void 0:_t.name}').update('RECORD_ID', body: { ... });
        batch.collection('${(yt=a[0])==null?void 0:yt.name}').delete('RECORD_ID');
        batch.collection('${(kt=a[0])==null?void 0:kt.name}').upsert(body: { ... });

        final result = await batch.send();
    `}}),R=new St({props:{language:"javascript",content:`
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
                        `}});let Y=Z(a[3]);const dt=e=>e[6].code;for(let e=0;e<Y.length;e+=1){let l=Et(a,Y,e),m=dt(l);nt.set(m,C[e]=Ot(m,l))}let N=Z(a[3]);const ut=e=>e[6].code;for(let e=0;e<N.length;e+=1){let l=Bt(a,N,e),m=ut(l);ct.set(m,$[e]=It(m,l))}return{c(){s=o("h3"),n=_("Batch create/update/upsert/delete ("),y=_(c),h=_(")"),r=i(),f=o("div"),f.innerHTML="<p>Batch and transactional create/update/upsert/delete of multiple records in a single request.</p>",F=i(),g=o("div"),U=o("div"),U.innerHTML='<i class="ri-error-warning-line"></i>',Fe=i(),P=o("div"),B=o("p"),Be=_(`The batch Web API need to be explicitly enabled and configured from the
            `),E=o("a"),E.textContent="Dashboard settings",Ee=_("."),Oe=i(),te=o("p"),te.textContent=`Because this endpoint process the requests in a single transaction, it could degrade the
            performance of your application if not used with proper care and configuration (e.g. too large
            allowed execution timeout, large body size limit, etc.).`,le=i(),Re(w.$$.fragment),oe=i(),O=o("h6"),O.textContent="API details",ae=i(),I=o("div"),I.innerHTML='<strong class="label label-primary">POST</strong> <div class="content">/api/batch</div>',se=i(),S=o("div"),S.textContent="Body Parameters",ne=i(),z=o("p"),z.innerHTML=`Body parameters could be sent as <em>application/json</em> or <em>multipart/form-data</em>.
    <br/>
    File upload is supported only via <em>multipart/form-data</em> (see below for more details).`,ie=i(),q=o("table"),ce=o("thead"),ce.innerHTML='<tr><th>Param</th> <th width="80%">Description</th></tr>',Ie=i(),re=o("tbody"),H=o("tr"),J=o("td"),J.innerHTML='<div class="flex txt-nowrap"><span class="label label-success">Required</span> <span>requests</span></div>',Se=i(),k=o("td"),W=o("span"),W.textContent="Array<Request>",He=_(` - List of the requests to process.

                `),de=o("p"),de.textContent="The supported batch request actions are:",Le=i(),D=o("ul"),G=o("li"),je=_("record create - "),ue=o("code"),ue.textContent="POST /api/collections/{collection}/records",Ae=i(),K=o("li"),Me=_(`record update -
                        `),pe=o("code"),pe.textContent="PATCH /api/collections/{collection}/records/{id}",Ne=i(),v=o("li"),Ue=_("record upsert - "),me=o("code"),me.textContent="PUT /api/collections/{collection}/records",ze=i(),Je=o("br"),We=i(),V=o("small"),V.innerHTML='(the body must have <code class="txt-sm">id</code> field)',Ge=i(),X=o("li"),Ke=_(`record delete -
                        `),be=o("code"),be.textContent="DELETE /api/collections/{collection}/records/{id}",Ve=i(),fe=o("p"),fe.textContent="Each batch Request element have the following properties:",Xe=i(),he=o("ul"),he.innerHTML=`<li><code>url path</code> <em>(could include query parameters)</em></li> <li><code>method</code> <em>(GET, POST, PUT, PATCH, DELETE)</em></li> <li><code>headers</code> <br/> <em>(if the <code>POST /api/batch</code> call has <code>Authorization</code> header it
                            is forwarded to all batch requests automatically)</em></li> <li><code>body</code></li>`,Ye=i(),p=o("p"),_e=o("strong"),_e.textContent="NB!",Ze=_(` When the batch request is send as
                    `),ye=o("code"),ye.textContent="multipart/form-data",Qe=_(`, the regular batch action fields are expected to be
                    submitted as serailized json under the `),ke=o("code"),ke.textContent="@jsonPayload",xe=_(` field and file keys need
                    to follow the pattern `),$e=o("code"),$e.textContent="requests.N.fileField",et=_(` or
                    `),ge=o("code"),ge.textContent="requests[N].fileField",tt=i(),Ce=o("em"),Ce.textContent=`(this is usually handled transparently by the SDKs when their specific object notation
                        is used)
                    `,lt=_(`.
                    `),ot=o("br"),at=_(`
                    If you don't use the SDKs or prefer manually to construct the `),De=o("code"),De.textContent="FormData",st=_(`
                    body, then it could look something like:
                    `),Re(R.$$.fragment),ve=i(),L=o("div"),L.textContent="Responses",we=i(),T=o("div"),j=o("div");for(let e=0;e<C.length;e+=1)C[e].c();it=i(),A=o("div");for(let e=0;e<$.length;e+=1)$[e].c();b(s,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(U,"class","icon"),b(E,"href","/settings"),b(P,"class","content"),b(g,"class","alert alert-warning"),b(O,"class","m-b-xs"),b(I,"class","api-route alert alert-success"),b(S,"class","section-title"),b(J,"valign","top"),b(W,"class","label"),b(V,"class","txt-hint"),b(q,"class","table-compact table-border m-t-xs m-b-base"),b(L,"class","section-title"),b(j,"class","tabs-header compact combined left"),b(A,"class","tabs-content"),b(T,"class","tabs")},m(e,l){d(e,s,l),t(s,n),t(s,y),t(s,h),d(e,r,l),d(e,f,l),d(e,F,l),d(e,g,l),t(g,U),t(g,Fe),t(g,P),t(P,B),t(B,Be),t(B,E),t(B,Ee),t(P,Oe),t(P,te),d(e,le,l),Te(w,e,l),d(e,oe,l),d(e,O,l),d(e,ae,l),d(e,I,l),d(e,se,l),d(e,S,l),d(e,ne,l),d(e,z,l),d(e,ie,l),d(e,q,l),t(q,ce),t(q,Ie),t(q,re),t(re,H),t(H,J),t(H,Se),t(H,k),t(k,W),t(k,He),t(k,de),t(k,Le),t(k,D),t(D,G),t(G,je),t(G,ue),t(D,Ae),t(D,K),t(K,Me),t(K,pe),t(D,Ne),t(D,v),t(v,Ue),t(v,me),t(v,ze),t(v,Je),t(v,We),t(v,V),t(D,Ge),t(D,X),t(X,Ke),t(X,be),t(k,Ve),t(k,fe),t(k,Xe),t(k,he),t(k,Ye),t(k,p),t(p,_e),t(p,Ze),t(p,ye),t(p,Qe),t(p,ke),t(p,xe),t(p,$e),t(p,et),t(p,ge),t(p,tt),t(p,Ce),t(p,lt),t(p,ot),t(p,at),t(p,De),t(p,st),Te(R,p,null),d(e,ve,l),d(e,L,l),d(e,we,l),d(e,T,l),t(T,j);for(let m=0;m<C.length;m+=1)C[m]&&C[m].m(j,null);t(T,it),t(T,A);for(let m=0;m<$.length;m+=1)$[m]&&$[m].m(A,null);M=!0,qe||(rt=Mt(Nt.call(null,E)),qe=!0)},p(e,[l]){var gt,Ct,Dt,vt,wt,qt,Rt,Tt;(!M||l&1)&&c!==(c=e[0].name+"")&&Ut(y,c);const m={};l&5&&(m.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[2]}');

        ...

        const batch = pb.createBatch();

        batch.collection('${(gt=e[0])==null?void 0:gt.name}').create({ ... });
        batch.collection('${(Ct=e[0])==null?void 0:Ct.name}').update("RECORD_ID", { ... });
        batch.collection('${(Dt=e[0])==null?void 0:Dt.name}').delete("RECORD_ID");
        batch.collection('${(vt=e[0])==null?void 0:vt.name}').upsert({ ... });

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
                        `),R.$set($t),l&10&&(Y=Z(e[3]),C=Pt(C,l,dt,1,e,Y,nt,j,zt,Ot,null,Et)),l&10&&(N=Z(e[3]),Jt(),$=Pt($,l,ut,1,e,N,ct,A,Wt,It,null,Bt),Gt())},i(e){if(!M){Q(w.$$.fragment,e),Q(R.$$.fragment,e);for(let l=0;l<N.length;l+=1)Q($[l]);M=!0}},o(e){x(w.$$.fragment,e),x(R.$$.fragment,e);for(let l=0;l<$.length;l+=1)x($[l]);M=!1},d(e){e&&(u(s),u(r),u(f),u(F),u(g),u(le),u(oe),u(O),u(ae),u(I),u(se),u(S),u(ne),u(z),u(ie),u(q),u(ve),u(L),u(we),u(T)),Pe(w,e),Pe(R);for(let l=0;l<C.length;l+=1)C[l].d();for(let l=0;l<$.length;l+=1)$[l].d();qe=!1,rt()}}}function Yt(a,s,n){let c,y,{collection:h}=s,r=200,f=[];const F=g=>n(1,r=g.code);return a.$$set=g=>{"collection"in g&&n(0,h=g.collection)},a.$$.update=()=>{a.$$.dirty&1&&n(4,y=Ft.dummyCollectionRecord(h)),a.$$.dirty&17&&h!=null&&h.id&&(f.push({code:200,body:JSON.stringify([y,Object.assign({},y,{id:y+"2"})],null,2)}),f.push({code:400,body:`
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
            `}),f.push({code:403,body:`
                {
                  "code": 403,
                  "message": "Batch requests are not allowed.",
                  "data": {}
                }
            `}))},n(2,c=Ft.getApiExampleUrl(Kt.baseURL)),[h,r,c,f,y,F]}class Qt extends Ht{constructor(s){super(),Lt(this,s,Yt,Xt,jt,{collection:0})}}export{Qt as default};
