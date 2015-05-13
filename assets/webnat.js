$( document ).ready(function () {

    var Template = function() {
        this.cached = {};
    };

    var T = new Template();
    $.extend(Template.prototype, {
        render: function(name, callback) {
            if (T.isCached(name)) {
                callback(T.cached[name]);
            } else {
                $.get(T.urlFor(name), function(raw) {
                    T.store(name, raw);
                    T.render(name, callback);
                });
            }
        },
        renderSync: function(name, callback) {
            if (!T.isCached(name)) {
                T.fetch(name);
            }
            T.render(name, callback);
        },
        prefetch: function(name) {
            $.get(T.urlFor(name), function(raw) { 
                T.store(name, raw);
            });
        },
        fetch: function(name) {
            if (! T.isCached(name)) {
                var raw = $.ajax({'url': T.urlFor(name), 'async': false}).responseText;
                T.store(name, raw);         
            }
        },
        isCached: function(name) {
            return !!T.cached[name];
        },
        store: function(name, raw) {
            T.cached[name] = Handlebars.compile(raw);
        },
        urlFor: function(name) {
            return "/assets/templates/"+ name + ".handlebars";
        }
    });

    // load data and then draw
    //
    async.parallel([
        function(callback) {
            $.getJSON("/shipyard_containers", function(data) {
                callback(null, data);
            });
        },
        function(callback) {
            $.getJSON("/ip_table", function(data) {
                callback(null, data); 
            });
        }],
    function (err, results) {
        context = {
            containers: results[0],
            ipTable: results[1]
        };

        drawGraph(context);
    });

    function drawGraph(graphContext) {
        
        var svg = d3.select("svg"),
            inner = svg.select("g"),
            zoom = d3.behavior.zoom().on("zoom", function() {
                inner.attr("transform", "translate(" + d3.event.translate + ")" +
                    "scale(" + d3.event.scale + ")");
            });
        svg.call(zoom);

        var render = new dagreD3.render();

        var g = new dagreD3.graphlib.Graph({
            compound: true,
            mulitgraph: false
        });
        g.setGraph({
            edgesep: 25,
            nodesep: 70,
            ranksep: 50,
            rankdir: "LR",
            marginx: 20,
            marginy: 20
        });

        var uniqSources = _.uniq(_.flatten(_.map(graphContext.ipTable, function(c) { return _.map(c.port_nat, function (p) { return p.from }) }) ));
        var sources = _.flatten(_.map(graphContext.ipTable, function(c) { return _.map(c.port_nat, function (p) { return p; }) }) );
        var hosts = _.uniq(_.map(graphContext.ipTable, function(c) { return c.host }));

        function htmlNode(caption, subtext, className) {
            var res;

            context = {
                caption: caption,
                subtext: subtext
            }
            T.renderSync('node', function(t) {
                res = t(context);
            });

            return res;
        }

        function draw(isUpdate) {

            // add NAT box
            var html = htmlNode("NAT", "NAT vm", "nat");
            console.log(html);
            g.setNode("nat-box", {
                labelType: "html",
                label: html,
                rx: 5,
                ry: 5,
                padding: 0
                //class: "nat-box"
            });

            for (var id in hosts) {
                var host = hosts[id];
                g.setNode(host, { 
                    label: host, 
                    labelType: "html",
                    rx: 5, 
                    ry: 5, 
                    width: 80 
                });
            }

            for (var id in uniqSources) {
                var source = uniqSources[id]
                    var html = htmlNode(source, source, "source");

                g.setNode("src_" + source, {
                    labelType: "html",
                    label: html,
                    rx: 5,
                    ry: 5,
                    padding: 0
                });

            }

            for (var id in graphContext.ipTable) {
                var container = graphContext.ipTable[id];
                var html = htmlNode(container.id, container.id, "container");

                g.setNode(id, {
                    labelType: "html",
                    label: html,
                    rx: 5,
                    ry: 5,
                    padding: 0
                    // class: className
                });

                g.setParent(id, container.host);

                if (container.port_nat) {
                    _.each(container.port_nat, function(nat) {
                        g.setEdge("nat-box", id, {
                            label: nat.port,
                            width: 40
                        });

                        g.setEdge("src_" + nat.from, "nat-box", {
                            label: nat.container_port,
                            width: 40
                        });
                    });
                }
            }

            inner.call(render, g);

            var zoomScale = zoom.scale();
            var graphWidth = g.graph().width + 80;
            var graphHeight = g.graph().height + 40;
            var width = parseInt(svg.style("width").replace(/px/, ""));
            var height = parseInt(svg.style("height").replace(/px/, ""));
            zoomScale = Math.min(width / graphWidth, height / graphHeight);
            var translate = [(width/2) - ((graphWidth*zoomScale)/2), (height/2) - ((graphHeight*zoomScale)/2)];
            zoom.translate(translate);
            zoom.scale(zoomScale);
            zoom.event(isUpdate ? svg.transition().duration(500) : d3.select("svg"));
        }

        draw(false);
    }
});
